
require('konva')

var stage;
var bglayer;
var fglayer;
var layer;
var msgtext;
var addedBeacons = [];
var group;

var origin = {};
origin.x = 307;
origin.y = 554;
var scale = {};
scale.y = -11.2;
scale.x = 10.9;

function writeMessage(message) {
  msgtext.setText(message);
  fglayer.draw();
}

$(function($) {
  setupKonva();
  populateBeacons();
});
  

function AddSelectedBeacon() {
    var selectedBeacon = $("#beaconlist").children('option').filter(":selected");
    var bid = selectedBeacon.val();
    var transform = stage.getAbsoluteTransform().copy();

            // to detect relative position we need to invert transform
            transform.invert();

            // now we find relative point
            var pos = stage.getPointerPosition();

            var mousepos = transform.point(pos);
    if(addBeacon(bid, mousepos.x, mousepos.y)) {
      selectedBeacon.attr('disabled', 'disabled')
      $("#beaconlist").val("")
    }
}

function addBeacon(bid, x, y) {

    if(bid=="") return false;

    var beacon = {};
    beacon.bid = bid;
    beacon.x = x;
    beacon.y = y;

    
    
    var label = new Konva.Label({
      x: x,
      y: y,
      draggable: true
    });

    label.add(new Konva.Tag({
            fill: 'black',
            pointerDirection: 'down',
            pointerWidth: 10,
            pointerHeight: 10,
            lineJoin: 'round',
            shadowColor: 'black',
            shadowBlur: 10,
            shadowOffset: 10,
            shadowOpacity: 0.5
        }));

    label.add(new Konva.Text({
        text: bid,
        fontFamily: 'Calibri',
        fontSize: 18,
        padding: 5,
        fill: 'white'
    }));

    label.on('dragmove', function () {
      var pos = this;
      var x = pos.x();
      var y = pos.y();
      writeMessage('x: ' + (x-origin.x)/scale.x + ', y: ' + (y-origin.y)/scale.y);
    });

    beacon.lbl = label;
    layer.add(label);
    layer.draw();

    addedBeacons.push(beacon);
    return true;
};

function populateBeacons() {
  $.get("/db/all/beacons", function(r){
      var beacons = r
      console.log(r)
      var select = $("#beaconlist")

      for(var b in beacons){
        var newoption = $("<option></option>")
        newoption.text(beacons[b])
        newoption.val(beacons[b])
        newoption.appendTo(select)
      }

      for(var b in beacons){
         $.get("/db/beacons/" + beacons[b] , function(r){
          console.log(r);
            addBeacon(r._id, r.xpos*scale.x + origin.x, r.ypos*scale.y + origin.y);
            var selectedBeacon = $("#beaconlist option[value="+r._id+"]");
            selectedBeacon.attr('disabled', 'disabled')
         });
      }
  });
}

function setupKonva() {
    // first we need to create a stage
  stage = new Konva.Stage({
    container: 'container',   // id of container <div>
    width: $('body').innerWidth()-10,
    height: $(window).innerHeight() - $('#headerbar').outerHeight(),
    draggable: true
  });
  bglayer = new Konva.Layer();
  fglayer = new Konva.Layer();
  // then create layer
  layer = new Konva.Layer();

  // create our shape
  circle = new Konva.Circle({
    x: stage.getWidth() / 2,
    y: stage.getHeight() / 2,
    radius: 70,
    fill: 'red',
    stroke: 'black',
    strokeWidth: 4,
    draggable: true
  });

  msgtext = new Konva.Text({
      x: 10,
      y: 10,
      fontFamily: 'Calibri',
      fontSize: 24,
      text: '',
      fill: 'black'
    });

  setFloorplan()

  stage.on('dblclick touchstart', AddSelectedBeacon);

  // add the shape to the layer
  layer.add(circle);
  //layer.draggable = true

  fglayer.add(msgtext)

  // add the layer to the stage
  stage.add(bglayer);
  stage.add(layer);
  stage.add(fglayer);

  // draw the image
  stage.draw();
}

function setFloorplan() {

  var imageObj = new Image()
  imageObj.onload = function() {
      var floorplan = new Konva.Image({
        x: 0,
        y: 0,
        image: this,
        width: this.width,
        height: this.height,
        //draggable: true
      });
      console.log(floorplan)
      bglayer.add(floorplan);
      bglayer.draw();
  };
  imageObj.src = "/db/buildings/eb2/floor1.png"
}


/*$("#floorplan").change(function(event) {
  var file = this.files[0]

  var imageObj = new Image()
  imageObj.onload = function() {
      var floorplan = new Konva.Image({
        x: 0,
        y: 0,
        image: this,
        width: this.width,
        height: this.height
      });
      console.log(floorplan)
      bglayer.add(floorplan);
      bglayer.draw();
  };
  imageObj.src = e.target.result

  var reader = new FileReader();
    reader.onload = function (e) {
      
    };
  reader.readAsDataURL(file);
});*/

$("#save").click(function() {
  for(var i = 0; i < addedBeacons.length; i++) {
    var beacon = addedBeacons[i];
    var x = beacon.lbl.x()
    var y = beacon.lbl.y()
    console.log({"xpos": (x-origin.x)/scale.x, "ypos": (y-origin.y)/scale.y});
      $.ajax({
        url: '/db/beacons/' + beacon.bid, 
        type: 'PUT',
        data: JSON.stringify({"xpos": (x-origin.x)/scale.x, "ypos": (y-origin.y)/scale.y})
      });
  }
});