import requests
import re
import json
import subprocess
import os
import pytz
import jdatetime
from urllib.parse import urlparse, parse_qs
import base64
import yaml

# لینک سابسکریپشن
SUBSCRIPTION_URL = "https://raw.githubusercontent.com/Giromo0/Collector/refs/heads/main/All_Configs_Sub.txt"

def fetch_configs():
    """دریافت کانفیگ‌ها از لینک سابسکریپشن"""
    try:
        response = requests.get(SUBSCRIPTION_URL, timeout=10)
        response.raise_for_status()
        return response.text.splitlines()
    except requests.RequestException as e:
        print(f"خطا در دریافت کانفیگ‌ها: {e}")
        return []

def parse_config(config):
    """پارس کردن کانفیگ و تبدیل به فرمت JSON برای Xray"""
    try:
        if config.startswith("vmess://"):
            # پارس vmess
            vmess_data = base64.b64decode(config.replace("vmess://", "")).decode()
            vmess_json = json.loads(vmess_data)
            return {
                "protocol": "vmess",
                "config": {
                    "v": "2",
                    "ps": vmess_json.get("ps", "unnamed"),
                    "add": vmess_json["add"],
                    "port": int(vmess_json["port"]),
                    "id": vmess_json["id"],
                    "aid": int(vmess_json.get("aid", 0)),
                    "net": vmess_json.get("net", "tcp"),
                    "type": vmess_json.get("type", "none"),
                    "host": vmess_json.get("host", ""),
                    "path": vmess_json.get("path", ""),
                    "tls": vmess_json.get("tls", "")
                }
            }
        elif config.startswith("vless://"):
            # پارس vless
            parsed = urlparse(config)
            params = parse_qs(parsed.query)
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
            # پارس ss
            parsed = urlparse(config)
            auth = base64.b64decode(parsed.netloc.split("@")[0]).decode().split(":")
            return {
                "protocol": "ss",
                "config": {
                    "method": auth[0],
                    "password": auth[1],
                    "add": parsed.hostname,
                    "port": int(parsed.port),
                    "ps": parsed.fragment or "unnamed"
                }
            }
        elif config.startswith("trojan://"):
            # پارس trojan
            parsed = urlparse(config)
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
        print(f"خطا در پارس کانفیگ: {e}")
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
    with open("temp_config.json", "w") as f:
        json.dump(xray_config, f)
    return "temp_config.json"

def test_config(config_file):
    """تست اتصال کانفیگ با Xray"""
    try:
        # اجرای Xray در پس‌زمینه
        process = subprocess.Popen(
            ["xray", "-c", config_file],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        # تست اتصال با curl
        result = subprocess.run(
            ["curl", "-x", "socks5://127.0.0.1:10808", "--connect-timeout", "5", "https://www.google.com"],
            capture_output=True,
            text=True
        )
        process.terminate()  # متوقف کردن Xray
        return result.returncode == 0
    except Exception as e:
        print(f"خطا در تست اتصال: {e}")
        return False

def extract_configs(lines):
    """استخراج و تست کانفیگ‌ها"""
    protocols = {"vless": [], "vmess": [], "ss": [], "trojan": []}
    pattern = r'^(vless://|vmess://|ss://|trojan://)[^\s#]+'

    for line in lines:
        match = re.match(pattern, line)
        if match:
            protocol = match.group(1).replace("://", "")
            if protocol in protocols and len(protocols[protocol]) < 20:
                parsed_config = parse_config(line)
                if parsed_config:
                    config_file = create_xray_config(parsed_config)
                    if test_config(config_file):
                        protocols[protocol].append(line.split("#")[0].strip())
    
    return protocols

def save_configs(protocols):
    """ذخیره کانفیگ‌های معتبر"""
    output_dir = "tested"
    os.makedirs(output_dir, exist_ok=True)
    output_file = os.path.join(output_dir, "config_test.txt")

    current_date_time = jdatetime.datetime.now(pytz.timezone('Asia/Tehran'))
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

def main():
    config_lines = fetch_configs()
    if not config_lines:
        print("هیچ کانفیگی دریافت نشد.")
        return

    protocols = extract_configs(config_lines)
    save_configs(protocols)
    print(f"کانفیگ‌های معتبر در tested/config_test.txt ذخیره شدند.")

if __name__ == "__main__":
    main()
