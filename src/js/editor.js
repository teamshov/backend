

$(document).ready(function(){
  var steps = $("#steps").steps();
 
// Add step
steps.steps("add", {
    title: "HTML code", 
    content: "<strong>HTML code</strong>"
});
});

var map = new ol.Map({
        target: 'map',
        layers: [
          new ol.layer.Tile({
            source: new ol.source.OSM()
          })
        ],
        view: new ol.View({
          center: ol.proj.fromLonLat([37.41, 8.82]),
          zoom: 4
        })
      });