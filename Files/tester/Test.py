import socket
import re
import os
import shutil
from datetime import datetime
import pytz
import jdatetime
import time
import random
import base64
import json
from concurrent.futures import ThreadPoolExecutor, as_completed

# مسیر پوشه پروتکل‌ها
PROTOCOL_DIR = "Splitted-By-Protocol"
# فایل‌های پروتکل
PROTOCOL_FILES = [
    "hysteria2.txt",
    "ss.txt",
    "ssr.txt",
    "trojan.txt",
    "tuic.txt",
    "vless.txt",
    "vmess.txt",
    "wireguard.txt"
]
# پوشه برای ذخیره نتایج
OUTPUT_DIR = "tested"
# فایل خروجی
OUTPUT_FILE = os.path.join(OUTPUT_DIR, "config_test.txt")
# حداکثر تعداد کانفیگ موفق برای هر پروتکل
MAX_SUCCESSFUL_CONFIGS = 20
# حداکثر تعداد کانفیگ برای تست
MAX_CONFIGS_TO_TEST = 300  # افزایش به 300 برای شانس بیشتر
# Timeout برای تست اتصال
TIMEOUT = 4  # افزایش به 4 ثانیه برای پروتکل‌های سنگین‌تر

# دیباگ: چاپ مسیر فعلی و چک کردن وجود پوشه پروتکل‌ها
print(f"Current working directory: {os.getcwd()}")
if os.path.exists(PROTOCOL_DIR):
    print(f"Protocol directory found: {PROTOCOL_DIR}")
else:
    print(f"Protocol directory NOT found: {PROTOCOL_DIR}")

# ایجاد پوشه نتایج اگر وجود نداشته باشه
if not os.path.exists(OUTPUT_DIR):
    os.makedirs(OUTPUT_DIR)
    print(f"Created output directory: {OUTPUT_DIR}")

# پاک کردن فایل‌های قدیمی در پوشه tested
if os.path.exists(OUTPUT_DIR):
    for file in os.listdir(OUTPUT_DIR):
        file_path = os.path.join(OUTPUT_DIR, file)
        if os.path.isfile(file_path):
            os.remove(file_path)
            print(f"Removed old file: {file_path}")

# تابع برای دیکد کردن کانفیگ‌های ss
def decode_ss_config(config):
    try:
        if config.startswith("ss://"):
            encoded = config.split("://")[1].split("#")[0]
            decoded = base64.urlsafe_b64decode(encoded + "=" * (4 - len(encoded) % 4)).decode()
            match = re.match(r"(.+?):(.+?)@(.+?):(\d+)", decoded)
            if match:
                return match.group(3), int(match.group(4))
    except Exception as e:
        print(f"Error decoding ss config: {config}, error: {str(e)}")
    return None, None

# تابع برای دیکد کردن کانفیگ‌های ssr
def decode_ssr_config(config):
    try:
        if config.startswith("ssr://"):
            encoded = config.split("://")[1].split("/")[0]
            decoded = base64.urlsafe_b64decode(encoded + "=" * (4 - len(encoded) % 4)).decode()
            parts = decoded.split(":")
            if len(parts) >= 6:
                host, port = parts[0], int(parts[1])
                return host, port
    except Exception as e:
        print(f"Error decoding ssr config: {config}, error: {str(e)}")
    return None, None

# تابع برای دیکد کردن کانفیگ‌های vmess
def decode_vmess_config(config):
    try:
        if config.startswith("vmess://"):
            encoded = config.split("://")[1].split("#")[0]
            decoded = base64.urlsafe_b64decode(encoded + "=" * (4 - len(encoded) % 4)).decode()
            data = json.loads(decoded)
            host = data.get("add")
            port = int(data.get("port"))
            return host, port
    except Exception as e:
        print(f"Error decoding vmess config: {config}, error: {str(e)}")
    return None, None

# تابع برای استخراج IP/دامنه و پورت از لینک پروتکل
def extract_host_port(config):
    # فیلتر کردن آدرس‌های نامعتبر
    def is_valid_host(host):
        if not host:
            return False
        # چک کردن IPv6 ناقص
        if host.startswith("[") and not host.endswith("]"):
            return False
        # چک کردن فرمت دامنه یا IP
        if "@" in host or not re.match(r"^[a-zA-Z0-9\.\-:]+$", host):
            return False
        # چک کردن طول دامنه
        if len(host) > 255:
            return False
        return True

    patterns = [
        r"(vless|trojan|hysteria2|tuic|wireguard)://.+?@(.+?):(\d+)",  # استاندارد
        r"(vless|trojan|hysteria2|tuic|wireguard)://(.+?):(\d+)"  # بدون uuid
    ]
    for pattern in patterns:
        match = re.match(pattern, config)
        if match:
            host = match.group(2)
            port = int(match.group(3))
            if is_valid_host(host):
                return host, port
    if config.startswith("ss://"):
        host, port = decode_ss_config(config)
        if is_valid_host(host):
            return host, port
    if config.startswith("ssr://"):
        host, port = decode_ssr_config(config)
        if is_valid_host(host):
            return host, port
    if config.startswith("vmess://"):
        host, port = decode_vmess_config(config)
        if is_valid_host(host):
            return host, port
    return None, None

# تابع تست TCP connection و محاسبه پینگ
def test_connection_and_ping(config, timeout=TIMEOUT):
    host, port = extract_host_port(config)
    if not host or not port:
        print(f"Invalid config format: {config}")
        return None
    try:
        start_time = time.time()
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        result = sock.connect_ex((host, port))
        sock.close()
        if result == 0:  # اتصال موفق
            ping_time = (time.time() - start_time) * 1000  # تبدیل به میلی‌ثانیه
            print(f"Successful connection to {host}:{port}, ping: {ping_time:.2f}ms")
            return {
                "config": config,
                "host": host,
                "port": port,
                "ping": ping_time
            }
        else:
            print(f"Failed connection to {host}:{port}, result: {result}")
        return None
    except (socket.gaierror, socket.timeout) as e:
        print(f"Error connecting to {host}:{port}: {str(e)}")
        return None

# تاریخ و زمان برای نام‌گذاری (جلیلی، تهران)
current_date_time = jdatetime.datetime.now(pytz.timezone('Asia/Tehran'))
current_month = current_date_time.strftime("%b")
current_day = current_date_time.strftime("%d")
updated_hour = current_date_time.strftime("%H")
updated_minute = current_date_time.strftime("%M")
final_string = f"{current_month}-{current_day} | {updated_hour}:{updated_minute}"

# لیست برای ذخیره تمام کانفیگ‌های موفق
all_successful_configs = []

try:
    # پردازش هر فایل پروتکل
    for protocol_file in PROTOCOL_FILES:
        file_path = os.path.join(PROTOCOL_DIR, protocol_file)
        protocol_name = protocol_file.replace(".txt", "").lower()
        
        # خواندن لینک‌های پروتکل از فایل
        config_links = []
        if os.path.exists(file_path):
            with open(file_path, 'r', encoding='utf-8') as f:
                config_links = [line.strip() for line in f if line.strip()]
                print(f"Found {len(config_links)} configs in {protocol_file}")
        else:
            print(f"Protocol file not found: {file_path}")
        
        # انتخاب تصادفی حداکثر 300 کانفیگ برای تست
        if len(config_links) > MAX_CONFIGS_TO_TEST:
            config_links = random.sample(config_links, MAX_CONFIGS_TO_TEST)
            print(f"Selected {len(config_links)} random configs for testing in {protocol_file}")
        
        # تست موازی کانفیگ‌ها
        configs_with_ping = []
        with ThreadPoolExecutor(max_workers=20) as executor:
            future_to_config = {executor.submit(test_connection_and_ping, config): config for config in config_links}
            for future in as_completed(future_to_config):
                try:
                    result = future.result()
                    if result and len(configs_with_ping) < MAX_SUCCESSFUL_CONFIGS:
                        result["protocol"] = protocol_name
                        configs_with_ping.append(result)
                except Exception as e:
                    print(f"Error processing config {future_to_config[future]}: {str(e)}")
        
        # مرتب‌سازی بر اساس پینگ و انتخاب حداکثر 20 کانفیگ
        configs_with_ping.sort(key=lambda x: x["ping"])
        successful_configs = configs_with_ping[:MAX_SUCCESSFUL_CONFIGS]
        print(f"Found {len(successful_configs)} successful configs for {protocol_name}")
        
        # اضافه کردن به لیست کلی
        all_successful_configs.extend(successful_configs)

except Exception as e:
    print(f"Error in main loop: {str(e)}")

# ذخیره تمام کانفیگ‌های موفق در یک فایل
try:
    with open(OUTPUT_FILE, "w", encoding="utf-8") as file:
        file.write(f"#🌐 به روزرسانی شده در {final_string} | MTSRVRS\n")
        if all_successful_configs:
            for i, result in enumerate(all_successful_configs, 1):
                config_string = f"#🌐سرور {i} | {result['protocol']} | {final_string} | Ping: {result['ping']:.2f}ms"
                file.write(f"{result['config']}{config_string}\n")
        else:
            file.write("# No successful configs found\n")
    print(f"Output file {OUTPUT_FILE} created")
    # دیباگ: چک کردن وجود فایل
    if os.path.exists(OUTPUT_FILE):
        print(f"Output file {OUTPUT_FILE} created successfully with size {os.path.getsize(OUTPUT_FILE)} bytes")
    else:
        print(f"Failed to create output file {OUTPUT_FILE}")
except Exception as e:
    print(f"Error writing output file: {str(e)}")
