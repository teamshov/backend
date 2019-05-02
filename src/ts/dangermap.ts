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

class ShovItemManager {
    constructor(public floor: Floor, public dburl: string, public color: string) {
        $.get('http://localhost:62027/db/all/' + dburl, (resp : any) => this.loadItems(resp));
    }

    loadItems(resp : any) {
        let ids: string[] = resp;
        for (let id of ids) {
            $.get('http://localhost:62027/db/' + this.dburl + '/' + id, (resp) => this.loadItem(resp));
        }
    }

    loadItem(doc : any) {
        if(this.floor.items[this.dburl] == null) {
            this.floor.items[this.dburl] = {}
        }
        this.floor.items[this.dburl][doc["_id"]] = doc;
        this.floor.items[this.dburl][doc["_id"]]["dangerlevel"] = 0.0
    }

}

class Cell {
    konvaObj : any;
    value : number;

    constructor(layer : any, public rpos: Vec2) {
        let kpos = toKPos(rpos);
        this.konvaObj = new Konva.Rect({
            x: kpos.x,
            y: kpos.y,
            width:15,
            height:15,
            fill: 'blue'
        })

        layer.add(this.konvaObj)
    }

    update(floor : Floor) {
        let v = 0;
        for(let dt in floor.items) {
            for(let d in floor.items[dt]) {
                let obj = floor.items[dt][d]
                if(obj["dangerlevel"] == null) {
                    continue;
                }
                let dpos = new Vec2(obj["xpos"], obj["ypos"])
                let dangerlevel = obj["dangerlevel"]
                v += this.probabilityFunc(this.rpos, dpos, dangerlevel);
            }
        }
        let p = 0;
        for(let pid in floor.positions) {
            let pos = floor.positions[pid]
            p += this.peopleProbabilityFunc(this.rpos, pos, 1);
        }
        this.konvaObj.fill(this.heatMapColorforValue(v, p))
    }

    probabilityFunc(upos:Vec2, bpos:Vec2, dangerlevel:number) {
        var dist = Math.sqrt((upos.x-bpos.x) ** 2 + (upos.y - bpos.y)**2);
        var A = Math.exp(-((dist)**2)/((2)**2));
        return A*dangerlevel;
    }

    peopleProbabilityFunc(upos:Vec2, bpos:Vec2, dangerlevel:number) {
        var dist = Math.sqrt((upos.x-bpos.x) ** 2 + (upos.y - bpos.y)**2);
        var A = Math.exp(-((dist)**2)/((0.8)**2));
        return A*dangerlevel;
    }
      
    heatMapColorforValue(dl : number, pl : number){
        var r = 255* (1-pl)
        var g = ((1-dl) * 255)* (1-pl);
        var b = ((1-dl) * 255 )
        return "rgb(" + r + ", "+g+","+ b+")";
    }
}

class HeatMap {

    updateInterval : any;
    grid : Cell[][];

    constructor(public floor : Floor) {
        this.grid = []
        for(let i = 0; i < 41; i++) {
            this.grid.push([])
            for(let j = 0; j < 57; j++) {
                this.grid[i].push(new Cell(floor.layer, new Vec2(i,j)))
            }
        }
        floor.layer.draw()

        this.updateInterval = setInterval(()=>{this.update()}, 1000);
    }

    update() {

        fetch("http://localhost:62027/api/data/eb2/L1").then(
            (d) => {
                d.json().then((data)=>{
                    this.updateData(data);
                })
            },
            (error) => {
                console.error(error)
            }
        ).catch((err)=>{
            console.error(err)
        })

        for(let i = 0; i < 41; i++) {
            for(let j = 0; j < 57; j++) {
                this.grid[i][j].update(this.floor);
            }
        }

        this.floor.layer.draw()
    }

    updateData(data:any) {
        this.floor.positions = []

        for(let k in data) {
            let keys = k.split(":")
            if(keys[0] == "device") {
                let deviceType =keys[1]
                let deviceID = keys[2]
                let dangerlevel = data[k]
                
                if(this.floor.items[deviceType] == null)
                    continue;
                if(this.floor.items[deviceType][deviceID] == null)
                    continue;

                this.floor.items[deviceType][deviceID]["dangerlevel"] = dangerlevel
            } else if(keys[0] == "users") {
                this.floor.positions.push(data[k]);
            }
        }
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

    items : any = {};

    positions : any = []

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
                opacity: 0.8
            });
            this.layer.add(this.floorplan);
            this.layer.draw();
        };
        this.imageObj.src = src

        this.esp32 = new ShovItemManager(this, "esp32", "blue");
        this.heatmap = new HeatMap(this);
        //this.graph = new ShovGraph(this, "http://localhost:62027/db/graphs/eb2_L1");
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
        this.floor = new Floor(this, "http://localhost:62027/db/buildings/eb2/L1_Black.png");
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
