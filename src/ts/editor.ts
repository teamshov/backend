import Konva from 'konva';
import $ from 'jquery';


class Vec2 {
    constructor(public x: number, public y: number) {}
}

class ELabel {
    label: any;
    text: any;
    tag: any;

    constructor(pos: Vec2, txt: string, fillColor: string) {
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


        this.label.add(this.tag);
        this.label.add(this.text);
    }
}

const scale: Vec2 = new Vec2(4.697624908, -4.697624908);
const origin: vec2 = new Vec2(0.32, 264.18);
const iscale: Vec2 = new Vec2(2, 2);

const toKPos = function (pos: Vec2): Vec2 {
    return new Vec2((pos.x * scale.x + origin.x)*iscale.x,(pos.y * scale.y + origin.y)*iscale.y);
}

const toPos = function (pos: Vec2): Vec2 {
    return new Vec2((pos.x/iscale.x - origin.x) / scale.x, (pos.y/iscale.y - origin.y) / scale.y);
}

class ShovItem {
    label: ELabel;
    constructor(public doc: any, color: string) {
        let pos = new Vec2(doc["xpos"], doc["ypos"]);
        this.label = new ELabel(toKPos(pos), doc["_id"], color);
        console.log(this.doc);
    }

    getKonvaObj() {
        return this.label.label;
    }

    save(url) {
        let kpos = this.label.label.position()
        let pos = toPos(new Vec2(kpos.x, kpos.y));
        $.ajax({
            url: url, 
            type: 'PUT',
            data: JSON.stringify({"xpos": pos.x, "ypos": pos.y})
          });
    }
}

class ShovItemManager {
    layer: any;
    ShovItems: ShovItem[] = [];

    constructor(public editor: Editor, public dburl: string, public color: string) {
        this.layer = new Konva.Layer({});
        editor.stage.add(this.layer);
        $.get('http://omaraa.ddns.net:62027/db/all/' + dburl, (resp) => this.loadItems(resp));
    }

    loadItems(resp) {
        let ids: string[] = resp;
        for (let id of ids) {
            $.get('http://omaraa.ddns.net:62027/db/' + this.dburl + '/' + id, (resp) => this.loadItem(resp));
        }
    }

    loadItem(doc) {
        let item = new ShovItem(doc, this.color);
        this.ShovItems.push(item);
        this.layer.add(item.getKonvaObj());
        this.layer.draw();
    }

    save() {
        for(let i of this.ShovItems) {
            i.save('http://omaraa.ddns.net:62027/db/' + this.dburl + '/' + i.doc['_id']);
        }
    }
}

class Building {
    layer: any;
    imageObj: any;
    floorplan: any;
    constructor(public editor: Editor, public src: string) {
        this.layer = new Konva.Layer({});
        editor.stage.add(this.layer);

        this.imageObj = new Image()
        this.imageObj.onload = () => {
            this.floorplan = new Konva.Image({
                x: 0,
                y: 0,
                image: this.imageObj,
                width: this.imageObj.width*iscale.x,
                height: this.imageObj.height*iscale.y,
                //draggable: true
            });
            this.layer.add(this.floorplan);
            this.layer.draw();
        };
        this.imageObj.src = src
    }


}

class Editor {
    stage: any;
    topbarElem: any;
    beacons: ShovItemManager;
    pies: ShovItemManger;
    constructor(public container: any, public topbarelem: any) {
        this.stage = new Konva.Stage({
            container: container, // id of container <div>
            width: container.offsetWidth,
            height: window.innerHeight - $('#topbar').height(),
            draggable: true,
            color: "white"
        });
        console.log("Konva initialized!");

        let layer = new Konva.Layer({});

        layer.add(new Konva.Circle({
            x: 0,
            y: 0,
            fill: "red",
            radius: 10
        }));

        this.stage.add(layer);
        this.bg = new Building(this, "http://omaraa.ddns.net:62027/db/buildings/eb2/L1.svg");
        this.beacons = new ShovItemManager(this, "beacons", "black");
        this.pies = new ShovItemManager(this, "pies", "red");

        this.stage.draw();

        window.save = () => {this.save();};
    }

    save() {
        this.beacons.save();
        this.pies.save();
    }
}

$(document).ready(() => {
    let container = document.getElementById('container');
    let topbar = document.getElementById('topbar');
    $("#topbar").ready(() => {
        let e = new Editor(container, topbar);
    })
})