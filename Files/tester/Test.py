import socket
import re
import os
import shutil
import ssl
import urllib.parse
from datetime import datetime
import pytz
import jdatetime
import time
import random
from concurrent.futures import ThreadPoolExecutor, as_completed
import websocket

# مسیر پوشه پروتکل‌ها
PROTOCOL_DIR = "Splitted-By-Protocol"

# لیست فایل‌های پروتکل (فقط پروتکل‌های موردنظر)
PROTOCOL_FILES = [
    "vmess.txt", "vless.txt", "trojan.txt", "ss.txt", "hy2.txt"
]

# پوشه برای ذخیره نتایج
OUTPUT_DIR = "tested"
# فایل خروجی
OUTPUT_FILE = os.path.join(OUTPUT_DIR, "config_test.txt")
# حداکثر تعداد کانفیگ موفق برای هر پروتکل
MAX_SUCCESSFUL_CONFIGS = 20
# حداکثر تعداد کانفیگ برای تست (برای پروتکل‌های اولویت‌دار)
MAX_CONFIGS_TO_TEST = 150  # افزایش برای پروتکل‌های موردنظر
# Timeout برای تست اتصال
TIMEOUT = float(os.getenv("TEST_TIMEOUT", 2))  # پیش‌فرض 2 ثانیه

# ایجاد پوشه نتایج اگر وجود نداشته باشه
if not os.path.exists(OUTPUT_DIR):
    os.makedirs(OUTPUT_DIR)

# پاک کردن فایل‌های قدیمی در پوشه tested
if os.path.exists(OUTPUT_DIR):
    for file in os.listdir(OUTPUT_DIR):
        file_path = os.path.join(OUTPUT_DIR, file)
        if os.path.isfile(file_path):
            os.remove(file_path)

# تابع برای استخراج IP/دامنه و پورت از لینک پروتکل
def extract_host_port(config):
    patterns = [
        r"(vmess|vless|trojan|ss|hy2)://.+?@(.+?):(\d+)",  # فرمت استاندارد
        r"(vmess|vless|trojan|ss|hy2)://(.+?):(\d+)"  # بدون uuid
    ]
    for pattern in patterns:
        match = re.match(pattern, config)
        if match:
            host = match.group(2)
            port = int(match.group(3))
            return host, port
    return None, None

# تابع برای پاکسازی توضیحات اضافی در کانفیگ
def clean_config(config, server_num, date_string):
    if "#" in config:
        main_config = config.split("#")[0]
        comment = urllib.parse.unquote(config.split("#")[1])
        # حذف کاراکترهای غیرالفبایی و غیرضروری
        cleaned_comment = "".join(c for c in comment if c.isalnum() or c in ".-_ ")
        return f"{main_config}#🌐 server-{server_num}-{date_string}"
    return f"{config}#🌐 server-{server_num}-{date_string}"

# تابع تست WebSocket
def test_websocket(config, timeout=TIMEOUT):
    host, port = extract_host_port(config)
    if not host or not port:
        return False
    try:
        parsed = urllib.parse.urlparse(config)
        path = parsed.query.split("path=")[1].split("&")[0] if "path=" in parsed.query else "/"
        ws_url = f"ws://{host}:{port}{path}"
        ws = websocket.create_connection(ws_url, timeout=timeout)
        ws.close()
        return True
    except Exception:
        return False

# تابع تست TLS
def test_tls(config, timeout=TIMEOUT):
    host, port = extract_host_port(config)
    if not host or not port:
        return False
    try:
        parsed = urllib.parse.urlparse(config)
        sni = parsed.query.split("sni=")[1].split("&")[0] if "sni=" in parsed.query else host
        context = ssl.create_default_context()
        with socket.create_connection((host, port), timeout=timeout) as sock:
            with context.wrap_socket(sock, server_hostname=sni) as ssock:
                return True
    except Exception:
        return False

# تابع تست TCP connection و محاسبه پینگ
def test_connection_and_ping(config, timeout=TIMEOUT):
    host, port = extract_host_port(config)
    if not host or not port:
        return None
    try:
        start_time = time.time()
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        result = sock.connect_ex((host, port))
        sock.close()
        if result == 0:  # اتصال TCP موفق
            ping_time = (time.time() - start_time) * 1000
            # تست WebSocket برای کانفیگ‌های ws
            if "type=ws" in config and not test_websocket(config, timeout):
                return None
            # تست TLS برای کانفیگ‌های tls
            if "security=tls" in config and not test_tls(config, timeout):
                return None
            return {
                "config": config,
                "host": host,
                "port": port,
                "ping": ping_time
            }
        return None
    except (socket.gaierror, socket.timeout, socket.error):
        return None

# تاریخ و زمان برای نام‌گذاری (جلیلی، تهران)
current_date_time = jdatetime.datetime.now(pytz.timezone('Asia/Tehran'))
current_month = current_date_time.strftime("%b")
current_day = current_date_time.strftime("%d")
updated_hour = current_date_time.strftime("%H")
updated_minute = current_date_time.strftime("%M")
final_string = f"{current_month}-{current_day}"

# لیست برای ذخیره تمام کانفیگ‌های موفق
all_successful_configs = []

# پردازش هر فایل پروتکل
for protocol_file in PROTOCOL_FILES:
    file_path = os.path.join(PROTOCOL_DIR, protocol_file)
    protocol_name = protocol_file.replace(".txt", "")
    
    # بررسی وجود فایل و غیرخالی بودن آن
    if not os.path.exists(file_path) or os.path.getsize(file_path) == 0:
        print(f"Skipping {file_path}: File is missing or empty")
        continue
    
    # خواندن لینک‌های پروتکل از فایل
    config_links = []
    with open(file_path, 'r', encoding='utf-8') as f:
        config_links = [line.strip() for line in f if line.strip() and not line.startswith('#')]
    
    # انتخاب تصادفی حداکثر 150 کانفیگ برای تست
    if len(config_links) > MAX_CONFIGS_TO_TEST:
        config_links = random.sample(config_links, MAX_CONFIGS_TO_TEST)
    
    # تست موازی کانفیگ‌ها
    configs_with_ping = []
    with ThreadPoolExecutor(max_workers=20) as executor:
        future_to_config = {executor.submit(test_connection_and_ping, config): config for config in config_links}
        for future in as_completed(future_to_config):
            result = future.result()
            if result and len(configs_with_ping) < MAX_SUCCESSFUL_CONFIGS:
                result["protocol"] = protocol_name
                configs_with_ping.append(result)
    
    # مرتب‌سازی بر اساس پینگ و انتخاب حداکثر 20 کانفیگ
    configs_with_ping.sort(key=lambda x: x["ping"])
    successful_configs = configs_with_ping[:MAX_SUCCESSFUL_CONFIGS]
    
    # اضافه کردن به لیست کلی
    all_successful_configs.extend(successful_configs)

# ذخیره تمام کانفیگ‌های موفق در فایل
if all_successful_configs:
    with open(OUTPUT_FILE, "w", encoding="utf-8") as file:
        file.write(f"# 🌐 به‌روزرسانی‌شده در {final_string} | MTSRVRS\n")
        for i, result in enumerate(all_successful_configs, 1):
            cleaned_config = clean_config(result["config"], i, final_string)
            config_string = f"# 🌐 server {i} | {result['protocol']} | {final_string} | Ping: {result['ping']:.2f}ms"
            file.write(f"{cleaned_config}\n{config_string}\n")
    print(f"All results saved to {OUTPUT_FILE}")
else:
    print("No successful configs found for any protocol")
