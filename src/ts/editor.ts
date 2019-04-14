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

declare global {
    interface Window { 
        save : any; 
        setNodeType : any;
    }
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
            draggable: true
        });


        this.label.add(this.tag);
        this.label.add(this.text);
    }
}

class NLine {
    line : any;
    constructor(public n1 : Node, public n2 : Node) {
        let p1 = this.n1.konvaObj.position();
        let p2 = this.n2.konvaObj.position();
        let arr : number[] = [p1.x, p1.y, p2.x, p2.y];
        this.line = new Konva.Line({
            points: arr,
            stroke: 'black',
            strokeWidth: 1,
        });

        this.n1.konvaObj.on("dragmove", () => this.onDragmove());
        this.n2.konvaObj.on("dragmove", () => this.onDragmove());

        n1.lines.add(this);
        n2.lines.add(this);
    }

    onDragmove() {
        if(this.n1.destroyed || this.n2.destroyed) {
            this.line.destroy();
            return;
        }
        let p1 = this.n1.konvaObj.position();
        let p2 = this.n2.konvaObj.position();
        let arr : number[] = [p1.x, p1.y, p2.x, p2.y];
        this.line.points(arr);
    }

    destroy(n : Node) {
        let othern = this.n1;
        if(n == this.n1) {
            othern == this.n2
        }

        othern.lines.delete(this);
        n.lines.delete(this);
        this.line.destroy();
    }
}

class Node {
    cNodes : Set<number> = new Set<number>();
    selected : boolean = false;
    konvaObj : any;
    dragged = false;
    destroyed = false;
    lines : Set<NLine> = new Set();
    fill : string;

    constructor(public graph : ShovGraph, pos : Vec2, public index : number, public type : string, nodes : number[] = null) {
        let kpos = toKPos(pos);
        this.fill = 'yellow';
        if(this.type == "exit") {
            this.fill = 'green';
        }
        this.konvaObj = new Konva.Circle({
            x: kpos.x,
            y: kpos.y,
            radius: 12,
            fill: this.fill,
            stroke: 'black',
            strokeWidth: 2,
            draggable: true,
        });


        this.konvaObj.on('mouseup touchend', () => this.onClick());
        this.konvaObj.on('dragstart', () => this.dragStart());
        this.konvaObj.on('dragend', () => this.dragEnd());

        if(nodes != null) {
            this.cNodes = new Set<number>(nodes);
        }
    }

    setPos(pos : Vec2) {
        this.konvaObj.position(pos);
    }

    getPos() : Vec2 {
        return this.konvaObj.position();
    }

    onClick() {
        if(!this.dragged) {
            this.graph.handleClick(this);
        }
    }

    dragStart() {
        this.dragged = true;
    }
    dragEnd() {
        this.dragged = false;
    }

    addNode(n : number) {
        if(!this.cNodes.has(n)) {
            this.cNodes.add(n);
            return true;
        } else {
            return false;
        }
    }

    toJSON() {
        return {"pos": toPos(this.getPos()), "cnodes": Array.from(this.cNodes), "index": this.index, "type": this.type};
    }
 }

class ShovGraph {
    layer : any;
    nodes : Map<number, Node> = new Map();
    nindex : number = 0;
    nodeType : string = '';
    dragged : boolean = false;

    selectednode : Node = null;

    constructor(public floor : Floor, src : string = null) {
        this.layer = new Konva.Layer({});
        floor.editor.stage.add(this.layer);
        
        if(src) {
            $.get(src, (resp : any) => this.loadGraph(resp));
        }

        floor.layer.on("mouseup touchend", ()=>{this.onStageClick()});
        floor.editor.stage.on('dragstart', () => this.dragStart());
        floor.editor.stage.on('dragend', () => this.dragEnd());
        this.layer.draw();

        window.setNodeType = (s:any) => {this.setNodeType(s)};
        let sinput = <HTMLInputElement>document.getElementById("nodeType");
        if(sinput)
        this.nodeType = sinput.value;
    }

    handleClick(node : Node) {
        if(this.selectednode == null){
            this.selectednode = node;
            node.konvaObj.fill("red");
        } else if(this.selectednode != node){
            this.connectNode(this.selectednode, node);
            this.selectednode.konvaObj.fill(this.selectednode.fill);
            this.selectednode = null;
        } else {
            this.selectednode.konvaObj.fill(this.selectednode.fill);
            this.selectednode = null;
            this.removeNode(node);
        }
        this.layer.draw();
    }
    dragStart() {
        this.dragged = true;
    }
    dragEnd() {
        this.dragged = false;
    }

    connectNode(n1 : Node, n2 : Node) {
        let r1 = n1.addNode(n2.index);
        let r2 = n2.addNode(n1.index);
        if(r1 && r2) {
            let line = new NLine(n1, n2);
            this.layer.add(line.line);
        }
    }

    loadGraph(resp : any) {
        let arr = resp["nodes"];
        for(let n of arr) {       
            let pos = n["pos"];
            let cnodes = n["cnodes"]
            let index = n["index"]
            let type = n["type"]
            if(index > this.nindex)
                this.nindex = index;
            this.addNode(pos, type, cnodes, index);
        }
        this.nindex++;

        let added : Set<number> = new Set<number>();
        for(let n of this.nodes.values()) {
            if(added.has(n.index)) { continue; }

            added.add(n.index);
            for(let nindex of Array.from(n.cNodes)) {
                if(added.has(nindex)) { continue; }

                let n2 = this.nodes.get(nindex);
                let line = new NLine(n, n2);
                this.layer.add(line.line);
            }
        }
    }

    save(url : any) {
        let json = '{"nodes":'+JSON.stringify(Array.from(this.nodes.values())) + '}';
        $.ajax({
            url: url, 
            type: 'PUT',
            data: json
          });
    }

    addNode(pos : Vec2, type : string, cnodes : number[] = null, index : number = null) {
        let i = index;
        if(!index) {
            i = this.nindex++;
        }
        let n = new Node(this, pos, i, type, cnodes);
        this.layer.add(n.konvaObj);
        this.nodes.set(n.index, n);
        this.layer.draw();
    }

    removeNode(n : Node) {
        let i = n.index;
        for(let nindex of n.cNodes) {
            let n2 = this.nodes.get(nindex);
            n2.cNodes.delete(i);
        }
        for(let line of n.lines) {
            line.destroy(n);
        }
        n.konvaObj.destroy();
        this.nodes.delete(n.index);
    }

    onStageClick() {
        if (!this.dragged) {

            if (this.selectednode) {
                this.selectednode.konvaObj.fill(this.selectednode.fill);
                this.selectednode.konvaObj.draw()
                this.selectednode = null;

                return;
            } else {
                let stage = this.floor.editor.stage;
                // what is transform of parent element?
                var transform = stage.getAbsoluteTransform().copy();

                // to detect relative position we need to invert transform
                transform.invert();

                // now we find relative point
                var pos = transform.point(stage.getPointerPosition());

                this.addNode(toPos(pos), this.nodeType);
                return;
            }
        }
    }

    setNodeType(selectObj : any) {
        console.log(selectObj)
        let val = selectObj.value;
        this.nodeType = val;
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
            data: JSON.stringify({"xpos": pos.x, "ypos": pos.y}),
            success: () => {console.log("save successful")},
            error: () => {console.error("save failed")}
          });
    }
}

class ShovItemManager {
    layer: any;
    ShovItems: ShovItem[] = [];

    constructor(public floor: Floor, public dburl: string, public color: string) {
        this.layer = new Konva.Layer({});
        floor.editor.stage.add(this.layer);
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
    }

    save() {
        for(let i of this.ShovItems) {
            i.save('http://omaraa.ddns.net:62027/db/' + this.dburl + '/' + i.doc['_id']);
        }
    }
}

class Floor {
    layer: any;
    imageObj: any;
    floorplan: any;
    beacons: ShovItemManager;
    pies: ShovItemManager;
    graph : ShovGraph;
    constructor(public editor: Editor, public src: string) {
        this.layer = new Konva.Layer({});
        editor.stage.add(this.layer);

        this.imageObj = new Image()
        this.imageObj.onload = () => {
            this.floorplan = new Konva.Image({
                x: 0,
                y: 0,
                image: this.imageObj,
                width: this.imageObj.width*iscale.x/11.25,
                height: this.imageObj.height*iscale.y/11.25,
                //draggable: true
            });
            this.layer.add(this.floorplan);
            this.layer.draw();
        };
        this.imageObj.src = src

        this.beacons = new ShovItemManager(this, "beacons", "black");
        this.pies = new ShovItemManager(this, "pies", "red");
        this.graph = new ShovGraph(this, "http://omaraa.ddns.net:62027/db/graphs/eb2_L1");
    }

    save() {
        this.beacons.save();
        this.pies.save();
        this.graph.save("http://omaraa.ddns.net:62027/db/graphs/eb2_L1");
    }


}

class Editor {
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

        layer.add(new Konva.Circle({
            x: 0,
            y: 0,
            radius: 8,
            fill: 'yellow',
            stroke: 'black',
            strokeWidth: 2
        }));

        this.stage.add(layer);
        this.floor = new Floor(this, "http://omaraa.ddns.net:62027/db/buildings/eb2/L1.png");
        

        this.stage.draw();

        window.save = () => {this.save();};
    }

    save() {
        this.floor.save();
    }
}

class MFEditor {
    editor : Editor;

    oncreate(vnode : any) {
        console.log(vnode.dom)
        this.editor = new Editor(vnode.dom)
    }

    view() {
        return m('div')
    }
}


FocusManager.showFocusOnlyOnTab();

let isDialogOpen = false;
let isDrawerOpen = false;
let selectedColor : any;


const Buttons = {
    view: () => {
      return m("[style=padding:5px;position:fixed;z-index:99]", [
          m(ButtonGroup, [
          m(Button, {
            iconLeft: Icons.SETTINGS,
            label: "",
            fluid: true,
            size: 'xl',
            onclick: () => (isDrawerOpen = true)
          }),
          m(Button, {
            iconLeft: Icons.SAVE,
            label: "Save",
            fluid: true,
            size: 'xl',
            onclick: ()=> {window.save()}
          })
        ])
      ]);
    }
  };

class MEditor {
  view(vnode : any) {
      const editor =  m(MFEditor);


    return m('span', [
        m(Buttons),
        editor
    ])
  }
};

m.mount(document.body, App);
