import subprocess
from pathlib import Path
import os
from datetime import datetime, timedelta
import time
import signal

TEMPLATE_CONFIG_PATH = "schedule_now.yaml"

SCHEDULE_TIME_DELAY = 4
CORN_PROCESS_TERMINATE = 40


def main():
    script_dir = Path(__file__).resolve().parent
    project_root = script_dir.parent
    build_corn_bin(str(project_root), get_corn_bin())

    template_content = process_template(
        os.path.join(script_dir, TEMPLATE_CONFIG_PATH), SCHEDULE_TIME_DELAY
    )

    config_path = get_config_path()
    curr_config_data = get_curr_config_data(config_path)

    with open(config_path, "w") as f:
        f.write(template_content)

    try:
        process = subprocess.Popen(
            ["corn", "run"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
        )

        try:
            time.sleep(CORN_PROCESS_TERMINATE)
            process.send_signal(signal.SIGINT)
            stdout, stderr = process.communicate(timeout=3)

        except subprocess.TimeoutExpired:
            process.kill()
            stdout, stderr = process.communicate()

        print("--- stdout ---")
        print(stdout)
        print()
        print("--- stderr ---")
        print(stderr)

    finally:
        if curr_config_data is not None:
            with open(config_path, "w") as f:
                f.write(curr_config_data)
        elif os.path.exists(config_path):
            os.remove(config_path)


def get_cron_str_later(later_secs: int):
    now = datetime.now()
    later = now + timedelta(seconds=later_secs)

    sec = later.second
    minute = later.minute
    hour = later.hour
    day_of_month = later.day
    month = later.month
    day_of_week = later.weekday()

    # In cron syntax Sunday is 0 or 7 (in python Monday = 0, Sunday = 6)
    cron_day_of_week = (day_of_week + 1) % 7

    cron_str = f"{sec} {minute} {hour} {day_of_month} {month} {cron_day_of_week}"
    return cron_str


def build_corn_bin(project_root: str, corn_bin_path: str):
    os.makedirs(os.path.dirname(corn_bin_path), exist_ok=True)
    build_command = ["go", "build", "-o", corn_bin_path, "main.go"]
    print(f"Building corn binary at {corn_bin_path}")
    subprocess.run(build_command, cwd=project_root, check=True)
    print("Build successful. Getting and replacing template content")


def get_config_path():
    return os.path.join(Path.home(), ".config/corn/corn.yaml")


def get_corn_bin():
    return os.path.join(Path.home(), ".local/bin/corn")


def process_template(template_path: str, template_time_delay_secs: int):
    template_content = ""
    cron_str_later = get_cron_str_later(template_time_delay_secs)
    with open(template_path, "r") as template_file:
        template_content = template_file.read()
    return template_content.replace("@TEST_TIME", cron_str_later)


def get_curr_config_data(config_path: str):
    config_data = None
    config_existed = os.path.exists(config_path)
    if config_existed:
        with open(config_path, "r") as f:
            config_data = f.read()
    return config_data


if __name__ == "__main__":
    main()
