# Master-Slave Multi-Chain Architecture: Environment Setup Guide

This project requires a robust containerized environment to run the Hyperledger Fabric Master-Slave topology. Instructions are provided for both Linux and Windows environments to ensure reproducibility.

## System Requirements
- **RAM:** 16GB Minimum (Highly Recommended for Docker/WSL2 stability)
- **Storage:** 50GB Free Space
- **Dependencies:** Docker, Docker Compose, Go (1.20+), Python (3.10+)

---

## Linux Setup (Ubuntu / Mint / Debian)

1. **System Update & Core Tools**
   ```bash
   sudo apt update && sudo apt upgrade -y
   sudo apt install -y curl git jq build-essential gcc make
   ```

2. **Install Docker Engine**
   ```bash
   sudo apt install -y docker.io docker-compose
   sudo systemctl enable --now docker
   sudo usermod -aG docker $USER
   # NOTE: Log out and log back in to apply Docker group permissions.

   ```


3. **Install Compilers & Python Tools**
   ```bash
   sudo apt install -y golang-go python3-pip python3-venv

   ```


4. **Provision Hyperledger Fabric Binaries**
   ```bash
   curl -sSLO [https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh](https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh)
   chmod +x install-fabric.sh
   ./install-fabric.sh docker samples binary

   ```



---

## Windows Setup (Windows 10 / 11 Pro)

1. **Enable Windows Subsystem for Linux (WSL2)**
Open PowerShell as Administrator and run:
   ```powershell
   wsl --install

   ```


*Restart your computer if prompted.*
2. **Install Docker Desktop**
* Download and install Docker Desktop for Windows.
* Open Docker Desktop Settings -> Resources -> WSL Integration. Ensure integration is turned on for your default WSL distro.


3. **Configure WSL Limits (Crucial for Memory Management)**
In your Windows user folder, create a `.wslconfig` file with the following limits to prevent Windows from freezing:
   ```ini
   [wsl2]
   memory=12GB
   processors=4
   swap=4GB

   ```


4. **Install Linux Tools inside WSL**
Open your WSL terminal (e.g., Ubuntu) and run the Linux setup commands exactly as written in the Linux Setup section above.
