brew install prometheus grafana


cat <<EOT >> /usr/local/etc/prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: "golang-app"
    static_configs:
      - targets: ["localhost:8080"]
EOT


brew services start prometheus


brew services start grafana

echo "Prometheus is running on http://localhost:9090"
echo "Grafana is running on http://localhost:3000"
