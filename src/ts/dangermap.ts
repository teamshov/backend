import 'construct-ui/lib/index.css'
import Konva from 'konva';
import $ from 'jquery';
import m from 'mithril';
import {
    Button,
    Icons,
    CustomSelect,
    ButtonGroup,
    Drawer,
    Dialog,
    SelectList,
    ListItem,
    FocusManager,
    Card,
    Icon,
    Grid,
    Col
  } from "construct-ui";

class Vec2 {
    constructor(public x: number, public y: number) {}
}

const scale: Vec2 = new Vec2(4.697624908, -4.697624908);
const origin: Vec2 = new Vec2(0.32, 264.18);
const iscale: Vec2 = new Vec2(4, 4);

const toKPos = function (pos: Vec2): Vec2 {
    return new Vec2((pos.x * scale.x + origin.x)*iscale.x,(pos.y * scale.y + origin.y)*iscale.y);
}

const toPos = function (pos: Vec2): Vec2 {
    return new Vec2((pos.x/iscale.x - origin.x) / scale.x, (pos.y/iscale.y - origin.y) / scale.y);

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
            shadowOpacity: 0.5
        });

        this.label = new Konva.Label({
            x: pos.x,
            y: pos.y,
        });


        this.label.add(this.tag);
        this.label.add(this.text);
    }
}

class ShovItem {
    label: ELabel;
    constructor(public doc: any, color: string) {
        let pos = new Vec2(doc["xpos"], doc["ypos"]);
        this.label = new ELabel(toKPos(pos), doc["_id"], color);
    }

    getKonvaObj() {
        return this.label.label;
    }

    save(url : any) {
        let kpos = this.label.label.position()
        let pos = toPos(new Vec2(kpos.x, kpos.y));
        $.ajax({
            url: url, 
            type: 'PUT',
            data: JSON.stringify({"xpos": pos.x, "ypos": pos.y, "building": "eb2", "floor": "L1"}),
            success: () => {console.log("save successful")},
            error: () => {console.error("save failed")}
          });
    }
}

class ShovItemManager {
    layer: any;
    ShovItems: ShovItem[] = [];
    shovDocs: any;

    constructor(public floor: Floor, public dburl: string, public color: string) {
        this.layer = new Konva.Layer({});
        floor.dangermap.stage.add(this.layer);
        $.get('http://omaraa.ddns.net:62027/db/all/' + dburl, (resp : any) => this.loadItems(resp));
    }

    loadItems(resp : any) {
        let ids: string[] = resp;
        for (let id of ids) {
            $.get('http://omaraa.ddns.net:62027/db/' + this.dburl + '/' + id, (resp) => this.loadItem(resp));
        }
    }

    loadItem(doc : any) {
        let item = new ShovItem(doc, this.color);
        this.ShovItems.push(item);
        this.layer.add(item.getKonvaObj());
        this.layer.draw();
        this.shovDocs[doc["_id"]] = doc
    }

}

class HeatMap {

    constructor(public floor : Floor) {

    }
}

class Floor {
    heatmaplayer : any;
    layer: any;
    imageObj: any;
    floorplan: any;
    beacons: ShovItemManager;
    pies: ShovItemManager;
    esp32: ShovItemManager;
    heatmap : HeatMap;
    constructor(public dangermap: DangerMap, public src: string) {
        this.layer = new Konva.Layer({});
        this.heatmaplayer = new Konva.Layer({});

        dangermap.stage.add(this.heatmaplayer);
        dangermap.stage.add(this.layer);

        this.imageObj = new Image()
        this.imageObj.onload = () => {
            this.floorplan = new Konva.Image({
                x: 0,
                y: 0,
                image: this.imageObj,
                width: this.imageObj.width*iscale.x/11.25,
                height: this.imageObj.height*iscale.y/11.25,
                //draggable: true
                opacity: 0.5
            });
            this.layer.add(this.floorplan);
            this.layer.draw();
        };
        this.imageObj.src = src

        this.beacons = new ShovItemManager(this, "beacons", "black");
        this.pies = new ShovItemManager(this, "pies", "red");
        this.esp32 = new ShovItemManager(this, "esp32", "blue");
        this.heatmap = new HeatMap(this);
        //this.graph = new ShovGraph(this, "http://omaraa.ddns.net:62027/db/graphs/eb2_L1");
    }

}

class DangerMap {
    stage: any;
    
    floor : Floor;
    constructor(public container: any) {
        this.stage = new Konva.Stage({
            container: container, // id of container <div>
            width: screen.availWidth,
            height: screen.availHeight,
            draggable: true,
            color: "white"
        });
        console.log("Konva initialized!");

        let layer = new Konva.Layer({});


        this.stage.add(layer);
        this.floor = new Floor(this, "http://omaraa.ddns.net:62027/db/buildings/eb2/L1_Black.png");
        this.stage.draw();
    }
}

class MDangerMap {
    editor : DangerMap;

    oncreate(vnode : any) {
        console.log(vnode.dom)
        this.editor = new DangerMap(vnode.dom)
    }

    view() {
        return m('div')
    }
}


FocusManager.showFocusOnlyOnTab();

class App {
  view(vnode : any) {
      const editor =  m(MDangerMap);


    return m('span', [
        editor
    ])
  }
};

m.mount(document.body, App);
