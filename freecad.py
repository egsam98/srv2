#!/usr/bin/env python
# -*- coding: utf-8 -*-

import json
import re
import Draft

TOP = 1
BOTTOM = -30
START = 1 + 2.5
END = 40.5 - 2.5
TRACK_Y = -3.96304547787

TABLE_WIDTH = 10
TABLE_PADDING_BOTTOM = 10
TABLE_PADDING_START = 0

DOC_NAME = "SAPR_Kurs"
INPUT_PATH = '/home/egor/go/src/srv2/out.json'
TEMPLATE_PATH = '/home/egor/go/src/srv2/A3L1 GOST.svg'


class Task:
    def __init__(self, task_id, p, e, periods):
        self.id = task_id
        self.p = p
        self.e = e
        self.periods = periods

    @staticmethod
    def decode(obj):
        return Task(str(obj["id"]), str(obj["p"]), str(obj["e"]), obj["periods"])

    @staticmethod
    def merge(tasks):
        _map = {}
        for task in tasks:
            if task.id in _map:
                _map[task.id].periods += task.periods
            else:
                _map[task.id] = task
        return _map.values()


def show():
    Gui.runCommand("Draft_Drawing")


def draw_table(tasks):
    Gui.activateWorkbench("DrawingWorkbench")

    start = START + TABLE_PADDING_START
    bottom = BOTTOM + TABLE_PADDING_BOTTOM
    block_len = float(TABLE_WIDTH) / 3
    block_height = 1.

    # head
    for j, head in enumerate(["id", "p", "e"]):
        pl = FreeCAD.Placement()
        pl.Base = FreeCAD.Vector(start + j * block_len, bottom, 0.0)
        Draft.makeRectangle(length=block_len, height=block_height, placement=pl, face=False, support=None)
        show()
        Draft.makeText(head, FreeCAD.Vector(start + j * block_len + 0.1, bottom + block_height / 4))
        show()

    # content
    for i, task in enumerate(tasks):
        i = -i - 1
        for j, value in zip(range(3), [task.id, task.p, task.e]):
            pl = FreeCAD.Placement()
            pl.Base = FreeCAD.Vector(start + j*block_len, bottom + i, 0.0)
            Draft.makeRectangle(length=block_len, height=block_height, placement=pl, face=False, support=None)
            show()
            Draft.makeText(value, FreeCAD.Vector(start + j*block_len + 0.1, bottom + i + block_height / 4))
            show()


def draw_track():
    Gui.activateWorkbench("DrawingWorkbench")

    points = [FreeCAD.Vector(x, TRACK_Y, 0.0) for x in [START, END]]
    Draft.makeWire(points, closed=False, face=True, support=None)
    show()
    for i in xrange(int(END-START)+1):
        points = [
            FreeCAD.Vector(START+i, TRACK_Y + 0.25, 0),
            FreeCAD.Vector(START+i, TRACK_Y - 0.25, 0)
        ]
        Draft.makeWire(points, closed=False, face=True, support=None)
        show()
        Draft.makeText(str(i), FreeCAD.Vector(START+i - 0.15, TRACK_Y - 0.7))
        show()


def draw_block(start, end, text):
    Gui.activateWorkbench("DrawingWorkbench")
    start = (START + start)
    end = (START + end)
    pl = FreeCAD.Placement()
    pl.Base = FreeCAD.Vector(start, TRACK_Y, 0.0)
    Draft.makeRectangle(length=end - start, height=2, placement=pl, face=False, support=None)
    show()
    Draft.makeText(text, FreeCAD.Vector(start + float(end - start)/2 - 0.15, TRACK_Y + 1))
    show()


# Create document
doc = App.newDocument(DOC_NAME)
App.ActiveDocument = doc

# Create drawing frame
doc.addObject('Drawing::FeaturePage', 'Page')
doc.Page.Template = TEMPLATE_PATH

# FreeCADGui.getDocument("SAPR_Kurs").getObject("Page").HintOffsetX = -50.00
# FreeCADGui.getDocument("SAPR_Kurs").getObject("Page").HintOffsetY = 0.00
# FreeCADGui.getDocument("SAPR_Kurs").getObject("Page").HintScale = 1.00

Gui.activateWorkbench("DraftWorkbench")

draw_track()
with open(INPUT_PATH) as f:
    tasks = [Task.decode(obj) for obj in json.load(f)]
    tasks = Task.merge(tasks)
    draw_table(tasks)
    for task in tasks:
        for period in task.periods:
            draw_block(float(period["start"])/1000, float(period["end"])/1000, task.id)

Gui.activateWorkbench("ArchWorkbench")

# execfile('/home/egor/go/src/srv2/freecad.py')
