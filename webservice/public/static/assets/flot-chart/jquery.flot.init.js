/**
* Theme: Montran Admin Template
* Author: Coderthemes
* Module/App: Flot-Chart
*/

var dlData = [];
var uplData = [];

!function($) {
    "use strict";

    var FlotChart = function() {
        this.$body = $("body")
        this.$realData = []
    };

    
    //creates Combine Chart
    FlotChart.prototype.createCombineGraph = function(selector, ticks, labels, datas) {
        
        var data = [{
            label: labels[0],
            data: datas[0],
            lines: {
                show: true,
                fill: true
            },
            points: {
                show: true
            }
        }, {
            label: labels[1],
            data: datas[1],
            lines: {
                show: true
            },
            points: {
                show: true
            }
        }, {
            label: labels[2],
            data: datas[2],
            bars: {
                show: true
            }
        }];
        var options = {
            xaxis: {
                ticks: ticks
            },
            series: {
                shadowSize: 0
            },
            grid: {
                hoverable: true,
                clickable: true,
                tickColor: "#f9f9f9",
                borderWidth: 1,
                borderColor: "#eeeeee"
            },
            colors: ["#33b86c", "#1a2942", "#60b1cc"],
            tooltip: true,
            tooltipOpts: {
                defaultTheme: false
            },
            legend: {
              position: 'nw'
            },
        };

        $.plot($(selector), data, options);
    },

        // graphs
        //initializing various charts and components
        FlotChart.prototype.init = function() {

          var combinelabels = ["DL", "UPL"];
          var combinedatas = [dlData, uplData];

          this.createCombineGraph("#combine-chart #combine-chart-container", null, combinelabels , combinedatas);
        },

    //init flotchart
    $.FlotChart = new FlotChart, $.FlotChart.Constructor = FlotChart
    
}(window.jQuery),

//initializing flotchart
function($) {
    "use strict";
        

         $.getJSON('/api/get?name=devicecheckstatus&deviceid=58aac5d076b6e50fc14c5a5a', function (data) {
            if (data['Result'] == 'OK') {
                var jsonData = data['Records'];

                var input = jsonData['input'];
                var output = jsonData['output'];

                var arr = [];

                for (var i = 0; i < input.length; i++) {
                    dlData[i] = [input[i][0], parseInt(input[i][1])] ;
                }
                for (var i = 0; i < output.length; i++) {
                    uplData[i] = [output[i][0], parseInt(output[i][1])] ;
                }
            }

           $.FlotChart.init()
        });


    
}(window.jQuery);



