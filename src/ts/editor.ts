import Konva from 'konva';
import $ from 'jquery';
import 'bootstrap';

interface Vec2 {
    x : number;
    y : number;
}

class ELabel {
    label : any;
    text : any;
    tag : any;

    constructor(pos:Vec2, txt:string, fillColor:string, private container : any) {
        this.text = new Konva.Text({
            text: txt,
            fontFamily: 'Calibri',
            fontSize: 18,
            padding: 5,
            fill: 'white'
        });

        this.tag = new Konva.Tag({
            fill: fillColor,
            pointerDirection: 'down',
            pointerWidth: 10,
            pointerHeight: 10,
            lineJoin: 'round',
            shadowColor: 'black',
            shadowBlur: 10,
            shadowOffset: 10,
            shadowOpacity: 0.5
        });

        this.label = new Konva.Label({
            x: pos.x,
            y: pos.y,
            draggable: true
          });
    }
}

class Editor {
    stage : any;
    constructor(containerID : string) {
        this.stage = new Konva.Stage({
            container: containerID,   // id of container <div>
            width: $('body').innerWidth()-10,
            height: $(window).innerHeight() - $('#headerbar').outerHeight(),
            draggable: true,
            color: "white"
        });

        let layer = new Konva.Layer({});

        layer.add(new Konva.Circle({
            x: 20,
            y: 20,
            fill: "red",
            radius: 10
        }));

        this.stage.add(layer);
        this.stage.draw();
    }
}

let e = new Editor("container");