#!/usr/bin/env bash
# =============================================================================
# Setup Monitoring: Prometheus + Grafana + Node Exporter
# =============================================================================
# Jalankan di VPS Ubuntu/Debian sebagai root atau dengan sudo.
#
# Cara pakai:
#   chmod +x setup-monitoring.sh
#   sudo ./setup-monitoring.sh
# =============================================================================

set -euo pipefail

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║   JValleyVerse — Monitoring Stack Installer                 ║"
echo "║   Prometheus + Grafana + Node Exporter                     ║"
echo "╚══════════════════════════════════════════════════════════════╝"

# ── Konfigurasi ──────────────────────────────────────────────────────────
PROMETHEUS_VERSION="2.53.0"
GRAFANA_VERSION="latest"
NODE_EXPORTER_VERSION="1.8.1"

PROMETHEUS_DIR="/etc/prometheus"
GRAFANA_DIR="/etc/grafana"
DATA_DIR="/var/lib/prometheus"

# ── 1. Update OS ─────────────────────────────────────────────────────────
echo ""
echo "📦 Updating system packages..."
apt-get update -qq && apt-get upgrade -y -qq

# ── 2. Install dependencies ──────────────────────────────────────────────
echo ""
echo "📦 Installing dependencies..."
apt-get install -y -qq wget curl tar gzip systemd

# ═══════════════════════════════════════════════════════════════════════════
# 3. PROMETHEUS
# ═══════════════════════════════════════════════════════════════════════════
echo ""
echo "📊 Installing Prometheus v${PROMETHEUS_VERSION}..."

# Buat user & group
id -u prometheus &>/dev/null || useradd --no-create-home --shell /bin/false prometheus

# Buat direktori
mkdir -p ${PROMETHEUS_DIR}
mkdir -p ${DATA_DIR}

# Download & extract
cd /tmp
wget -q "https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz"
tar xzf "prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz"
cd "prometheus-${PROMETHEUS_VERSION}.linux-amd64"

# Copy binary
cp prometheus promtool /usr/local/bin/
chown prometheus:prometheus /usr/local/bin/prometheus /usr/local/bin/promtool

# Copy console libraries (opsional)
cp -r consoles console_libraries ${PROMETHEUS_DIR}/
chown -R prometheus:prometheus ${PROMETHEUS_DIR} ${DATA_DIR}

# Copy prometheus.yml (custom config — pastikan file sudah ada di repo!)
if [ -f "${PROMETHEUS_DIR}/prometheus.yml" ]; then
    echo "✅ prometheus.yml already exists at ${PROMETHEUS_DIR}"
else
    echo "⚠️  Please deploy deploy/monitoring/prometheus.yml to ${PROMETHEUS_DIR}/prometheus.yml"
    echo "   Example: cp deploy/monitoring/prometheus.yml ${PROMETHEUS_DIR}/prometheus.yml"
fi

# Systemd service
cat > /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=Prometheus Monitoring
Wants=network-online.target
After=network-online.target

[Service]
User=prometheus
Group=prometheus
Type=simple
ExecStart=/usr/local/bin/prometheus \\
    --config.file=${PROMETHEUS_DIR}/prometheus.yml \\
    --storage.tsdb.path=${DATA_DIR} \\
    --web.console.templates=${PROMETHEUS_DIR}/consoles \\
    --web.console.libraries=${PROMETHEUS_DIR}/console_libraries \\
    --web.listen-address=0.0.0.0:9090

Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable prometheus
systemctl restart prometheus

echo "✅ Prometheus running on http://localhost:9090"

# ═══════════════════════════════════════════════════════════════════════════
# 4. NODE EXPORTER
# ═══════════════════════════════════════════════════════════════════════════
echo ""
echo "💻 Installing Node Exporter v${NODE_EXPORTER_VERSION}..."

id -u node_exporter &>/dev/null || useradd --no-create-home --shell /bin/false node_exporter

cd /tmp
wget -q "https://github.com/prometheus/node_exporter/releases/download/v${NODE_EXPORTER_VERSION}/node_exporter-${NODE_EXPORTER_VERSION}.linux-amd64.tar.gz"
tar xzf "node_exporter-${NODE_EXPORTER_VERSION}.linux-amd64.tar.gz"
cd "node_exporter-${NODE_EXPORTER_VERSION}.linux-amd64"

cp node_exporter /usr/local/bin/
chown node_exporter:node_exporter /usr/local/bin/node_exporter

# Systemd service
cat > /etc/systemd/system/node_exporter.service <<EOF
[Unit]
Description=Node Exporter (System Metrics)
Wants=network-online.target
After=network-online.target

[Service]
User=node_exporter
Group=node_exporter
Type=simple
ExecStart=/usr/local/bin/node_exporter \\
    --web.listen-address=:9100 \\
    --collector.logind

Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable node_exporter
systemctl restart node_exporter

echo "✅ Node Exporter running on http://localhost:9100/metrics"

# ═══════════════════════════════════════════════════════════════════════════
# 5. GRAFANA
# ═══════════════════════════════════════════════════════════════════════════
echo ""
echo "📈 Installing Grafana (${GRAFANA_VERSION})..."

# Install Grafana dari repo resmi
apt-get install -y -qq software-properties-common
wget -q -O /usr/share/keyrings/grafana.key https://apt.grafana.com/gpg.key
echo "deb [signed-by=/usr/share/keyrings/grafana.key] https://apt.grafana.com stable main" \
    | tee /etc/apt/sources.list.d/grafana.list

apt-get update -qq
apt-get install -y -qq grafana

# Systemd service (already created by package)
systemctl enable grafana-server
systemctl restart grafana-server

echo ""
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║   ✅ INSTALASI SELESAI!                                    ║"
echo "║                                                            ║"
echo "║   Prometheus:  http://<VPS_IP>:9090                       ║"
echo "║   Grafana:     http://<VPS_IP>:3000  (admin / admin)      ║"
echo "║   Node Export: http://<VPS_IP>:9100/metrics               ║"
echo "║                                                            ║"
echo "║   📋 Langkah selanjutnya:                                 ║"
echo "║   1. Buka Grafana → Add data source → Prometheus          ║"
echo "║      URL: http://localhost:9090                           ║"
echo "║   2. Import dashboard → upload grafana-dashboard.json     ║"
echo "║   3. Ganti password Grafana default!                      ║"
echo "╚══════════════════════════════════════════════════════════════╝"
