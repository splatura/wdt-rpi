package main

import (
        "log"
        "net/http"
        "os/exec"
        "strconv"
        "time"
        "strings"

        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
        externalValue = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "wdt_input_voltage",
                Help: "Value returned by wdt -g v",
        })
)

var (
        externalValue2 = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: "wdt_battery_voltage",
                Help: "Value returned by wdt -g vb",
        })
)

func fetchExternalValue(param string) (int, error) {
        cmd := exec.Command("wdt", "-g", param)
        output, err := cmd.Output()
        if err != nil {
                return 0, err
        }
        valueStr := strings.TrimSpace(string(output)) // Trim leading/trailing whitespace
        value, err := strconv.Atoi(string(valueStr))
        if err != nil {
                return 0, err
        }
        return value, nil
}

func updateMetric() {
        value, err := fetchExternalValue("v")
        if err != nil {
                log.Printf("Failed to fetch external value: %v", err)
                return
        }
        externalValue.Set(float64(value))

        value2, err2 := fetchExternalValue("vb")
        if err2 != nil {
                log.Printf("Failed to fetch external value: %v", err2)
                return
        }
        externalValue2.Set(float64(value2))
}

func main() {
        prometheus.MustRegister(externalValue)
        prometheus.MustRegister(externalValue2)

        http.Handle("/metrics", promhttp.Handler())

        go func() {
                for {
                        updateMetric()
                        time.Sleep(1 * time.Minute)
                }
        }()

        log.Fatal(http.ListenAndServe(":9101", nil))
}
