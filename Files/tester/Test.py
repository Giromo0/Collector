import requests
import re
import json
import subprocess
import os
import pytz
import jdatetime
from urllib.parse import urlparse, parse_qs
import base64
from concurrent.futures import ThreadPoolExecutor
import logging
import uuid

# تنظیم لاگ‌گذاری
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

# لینک سابسکریپشن
SUBSCRIPTION_URL = "https://raw.githubusercontent.com/Giromo0/Collector/refs/heads/main/All_Configs_Sub.txt"

def fetch_configs():
    """دریافت کانفیگ‌ها از لینک سابسکریپشن"""
    try:
        response = requests.get(SUBSCRIPTION_URL, timeout=10)
        response.raise_for_status()
        return response.text.splitlines()
    except requests.RequestException as e:
        logging.error(f"خطا در دریافت کانفیگ‌ها: {e}")
        return []

def parse_config(config):
    """پارس کردن کانفیگ و تبدیل به فرمت JSON برای Xray"""
    try:
        if config.startswith("vmess://"):
            try:
                vmess_data = base64.b64decode(config.replace("vmess://", "") + "===").decode()
                vmess_json = json.loads(vmess_data)
                return {
                    "protocol": "vmess",
                    "config": {
                        "v": "2",
                        "ps": vmess_json.get("ps", "unnamed"),
                        "add": vmess_json.get("add", ""),
                        "port": int(vmess_json.get("port", 0)) or 443,
                        "id": vmess_json.get("id", ""),
                        "aid": int(vmess_json.get("aid", 0)),
                        "net": vmess_json.get("net", "tcp"),
                        "type": vmess_json.get("type", "none"),
                        "host": vmess_json.get("host", ""),
                        "path": vmess_json.get("path", ""),
                        "tls": vmess_json.get("tls", "")
                    }
                }
            except (base64.binascii.Error, json.JSONDecodeError) as e:
                logging.error(f"خطا در پارس vmess: {e} - کانفیگ: {config[:50]}...")
                return None
        elif config.startswith("vless://"):
            parsed = urlparse(config)
            params = parse_qs(parsed.query)
            if not parsed.hostname or not parsed.username:
                logging.error(f"کانفیگ vless ناقص است: hostname یا username وجود ندارد - کانفیگ: {config[:50]}...")
                return None
            return {
                "protocol": "vless",
                "config": {
                    "id": parsed.username,
                    "add": parsed.hostname,
                    "port": int(parsed.port or 443),
                    "net": params.get("type", ["tcp"])[0],
                    "path": params.get("path", [""])[0],
                    "security": params.get("security", ["none"])[0],
                    "ps": params.get("ps", ["unnamed"])[0]
                }
            }
        elif config.startswith("ss://"):
            try:
                parsed = urlparse(config)
                auth = base64.b64decode(parsed.netloc.split("@")[0] + "===").decode().split(":")
                return {
                    "protocol": "ss",
                    "config": {
                        "method": auth[0],
                        "password": auth[1],
                        "add": parsed.hostname,
                        "port": int(parsed.port or 443),
                        "ps": parsed.fragment or "unnamed"
                    }
                }
            except (base64.binascii.Error, IndexError) as e:
                logging.error(f"خطا در پارس ss: {e} - کانفیگ: {config[:50]}...")
                return None
        elif config.startswith("trojan://"):
            parsed = urlparse(config)
            if not parsed.hostname or not parsed.username:
                logging.error(f"کانفیگ trojan ناقص است: hostname یا username وجود ندارد - کانفیگ: {config[:50]}...")
                return None
            return {
                "protocol": "trojan",
                "config": {
                    "password": parsed.username,
                    "add": parsed.hostname,
                    "port": int(parsed.port or 443),
                    "ps": parsed.fragment or "unnamed"
                }
            }
        return None
    except Exception as e:
        logging.error(f"خطا در پارس کانفیگ: {e} - کانفیگ: {config[:50]}...")
        return None

def create_xray_config(parsed_config):
    """ایجاد فایل JSON برای Xray"""
    xray_config = {
        "inbounds": [
            {
                "port": 10808,
                "protocol": "socks",
                "settings": {"auth": "noauth", "udp": True}
            }
        ],
        "outbounds": [
            {
                "protocol": parsed_config["protocol"],
                "settings": {
                    parsed_config["protocol"]: [parsed_config["config"]]
                }
            },
            {"protocol": "freedom", "tag": "direct"}
        ],
        "routing": {
            "rules": [
                {"type": "field", "outboundTag": "direct", "domain": ["geosite:category-ads-all"]}
            ]
        }
    }
    config_file = f"temp_config_{uuid.uuid4().hex}.json"
    with open(config_file, "w") as f:
        json.dump(xray_config, f)
    return config_file

def test_config(config):
    """تست اتصال کانفیگ با Xray"""
    parsed_config = parse_config(config)
    if not parsed_config:
        return False, config

    config_file = create_xray_config(parsed_config)
    try:
        process = subprocess.Popen(
            ["xray", "-c", config_file],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        result = subprocess.run(
            ["curl", "-x", "socks5://127.0.0.1:10808", "--connect-timeout", "3", "http://1.1.1.1"],
            capture_output=True,
            text=True,
            timeout=5
        )
        stdout, stderr = process.communicate()
        process.terminate()
        try:
            os.remove(config_file)
        except FileNotFoundError:
            logging.warning(f"فایل {config_file} قبلاً حذف شده یا وجود ندارد")
        if result.returncode == 0:
            logging.info(f"کانفیگ با موفقیت تست شد: {config[:50]}...")
            return True, config
        else:
            logging.warning(f"تست کانفیگ شکست خورد: {result.stderr} - کانفیگ: {config[:50]}...")
            return False, config
    except Exception as e:
        logging.error(f"خطا در تست اتصال: {e} - کانفیگ: {config[:50]}...")
        try:
            os.remove(config_file)
        except FileNotFoundError:
            logging.warning(f"فایل {config_file} قبلاً حذف شده یا وجود ندارد")
        return False, config

def extract_configs(lines):
    """استخراج و تست کانفیگ‌ها"""
    protocols = {"vless": [], "vmess": [], "ss": [], "trojan": []}
    pattern = r'^(vless://|vmess://|ss://|trojan://)[^\s#]+'
    invalid_configs = []

    valid_configs = []
    with ThreadPoolExecutor(max_workers=5) as executor:
        configs = [line for line in lines if re.match(pattern, line)]
        logging.info(f"تعداد کانفیگ‌های یافت‌شده: {len(configs)}")
        results = executor.map(test_config, configs[:50])  # محدود به ۵۰ کانفیگ
        for is_valid, config in results:
            if is_valid:
                protocol = config.split("://")[0]
                if protocol in protocols and len(protocols[protocol]) < 20:
                    protocols[protocol].append(config.split("#")[0].strip())
            else:
                invalid_configs.append(config)
    
    # ذخیره کانفیگ‌های نامعتبر
    if invalid_configs:
        output_dir = "tested"
        os.makedirs(output_dir, exist_ok=True)
        with open(os.path.join(output_dir, "invalid_configs.txt"), "w", encoding="utf-8") as f:
            for config in invalid_configs:
                f.write(f"{config}\n")
        logging.info(f"کانفیگ‌های نامعتبر در tested/invalid_configs.txt ذخیره شدند.")

    return protocols

def save_configs(protocols):
    """ذخیره کانفیگ‌های معتبر"""
    output_dir = "tested"
    os.makedirs(output_dir, exist_ok=True)
    current_date_time = jdatetime.datetime.now(pytz.timezone('Asia/Tehran'))
    output_file = os.path.join(output_dir, f"config_test_{current_date_time.strftime('%Y-%m-%d_%H-%M')}.txt")

    final_string = current_date_time.strftime("%b-%d | %H:%M")
    final_others_string = current_date_time.strftime("%b-%d")

    with open(output_file, "w", encoding="utf-8") as file:
        i = 0
        for protocol, configs in protocols.items():
            for config in configs:
                config_string = (
                    f"#🌐 به روزرسانی شده در {final_string} | هر 15 دقیقه کانفیگ جدید داریم"
                    if i == 0
                    else f"#🌐سرور {i} | {final_others_string} | {protocol.upper()}"
                )
                file.write(f"{config}{config_string}\n")
                i += 1
    logging.info(f"کانفیگ‌های معتبر در {output_file} ذخیره شدند.")

def main():
    config_lines = fetch_configs()
    if not config_lines:
        logging.error("هیچ کانفیگی دریافت نشد.")
        return

    protocols = extract_configs(config_lines)
    save_configs(protocols)

if __name__ == "__main__":
    main()
