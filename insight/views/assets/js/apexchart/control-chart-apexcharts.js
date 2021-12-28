'use strict';
$(function() {
    chartD();
});

/* Vertical Bar Chart */

function chartD() {

    var options = {
        chart: {
            height: 350,
            type: 'bar',
            fontFamily: 'Poppins, sans-serif',
            toolbar: {
                show: false
            },
            zoom: {
                enabled: false
            },
        },

        plotOptions: {
            bar: {
                dataLabels: {
                    position: 'top', // top, center, bottom
                },
            }
        },
        colors: ["#11a0fd"],
        dataLabels: {
            enabled: true,
            formatter: function(val) {
                return val + "%";
            },
            offsetY: -20,
            style: {
                fontSize: '12px',
                colors: ["#10163a"],
                fontFamily: 'Poppins, sans-serif',
            }
        },
        series: [{
            name: 'Revenue',
            data: [2.3, 3.1, 4.0, 10.1, 4.0, 3.6, 3.2, 2.3, 1.4, 0.8, 0.5, 0.2]
        }],
        xaxis: {
            categories: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
            position: 'top',
            labels: {
                offsetY: -18,
                style: {
                    colors: '#10163a',
                    fontFamily: 'Poppins, sans-serif',
                }
            },
            axisBorder: {
                show: false
            },
            axisTicks: {
                show: false
            },
            crosshairs: {
                fill: {
                    type: 'gradient',
                    gradient: {
                        colorFrom: '#11a0fd',
                        colorTo: '#1b4962',
                        stops: [0, 100],
                        opacityFrom: 0.4,
                        opacityTo: 0.5,
                    }
                }
            },
            tooltip: {
                enabled: true,
                offsetY: -35,

            }
        },
        fill: {
            gradient: {
                shade: 'light',
                type: "horizontal",
                shadeIntensity: 0.25,
                gradientToColors: undefined,
                inverseColors: true,
                opacityFrom: 1,
                opacityTo: 1,
                stops: [50, 0, 100, 100]
            },
        },
        yaxis: {
            axisBorder: {
                show: false
            },
            axisTicks: {
                show: false,
            },
            labels: {
                show: false,
                formatter: function(val) {
                    return val + "%";
                }
            }

        },
        title: {
            text: 'Monthly Revenue',
            floating: true,
            offsetY: 320,
            align: 'center',
            style: {
                color: '#10163a',
                fontFamily: 'Poppins, sans-serif',
            }
        },
    }

    var chart = new ApexCharts(
        document.querySelector("#chartD"),
        options
    );

    chart.render();

}
/* Column Bar Chart */
