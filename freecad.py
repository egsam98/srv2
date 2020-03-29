import json
import Draft

TOP = 1
BOTTOM = -30
START = 1 + 2.5
END = 40.5 - 2.0
TRACK_Y = -10
TRACK_MARGIN = 5

TABLE_WIDTH = 20
TABLE_PADDING_BOTTOM = 10
TABLE_PADDING_START = -2

FRAME = 3
U = 0.741

DOC_NAME = "SAPR_Kurs"

# Пришлось...
INPUT_PATH = '/home/egor/go/src/srv2/out.json'
TEMPLATE_PATH = '/home/egor/go/src/srv2/A3L1 GOST.svg'


class Task:
    def __init__(self, task_id, p, e, periods):
        self.id = task_id
        self.p = p
        self.e = e
        self.periods = periods

    def u(self):
        return str(round(float(self.e) / float(self.p), 3))

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
        return list(_map.values())


def show():
    Gui.runCommand("Draft_Drawing")


def draw_table(tasks):
    Gui.activateWorkbench("DrawingWorkbench")

    tasks.sort(key=lambda t: t.id)

    start = START + TABLE_PADDING_START
    bottom = BOTTOM + TABLE_PADDING_BOTTOM
    block_len = float(TABLE_WIDTH) / 4
    block_height = 1.

    # head
    for j, head in enumerate(["id", "p", "e", "u_i (итого: {})".format(U)]):
        pl = FreeCAD.Placement()
        pl.Base = FreeCAD.Vector(start + j * block_len, bottom, 0.0)
        Draft.makeRectangle(length=block_len, height=block_height, placement=pl, face=False, support=None)
        show()
        Draft.makeText(head, FreeCAD.Vector(start + j * block_len + 0.1, bottom + block_height / 4))
        show()

    # content
    for i, task in enumerate(tasks):
        i = -i - 1
        for j, value in enumerate([task.id, task.p, task.e, task.u()]):
            pl = FreeCAD.Placement()
            pl.Base = FreeCAD.Vector(start + j*block_len, bottom + i, 0.0)
            Draft.makeRectangle(length=block_len, height=block_height, placement=pl, face=False, support=None)
            show()
            Draft.makeText(value, FreeCAD.Vector(start + j*block_len + 0.1, bottom + i + block_height / 4))
            show()


def draw_track(y, start_num):
    Gui.activateWorkbench("DrawingWorkbench")

    points = [FreeCAD.Vector(x, y, 0.0) for x in [START, END]]
    Draft.makeWire(points, closed=False, face=True, support=None)
    show()
    for i in range(int(END-START)+1):
        points = [
            FreeCAD.Vector(START+i, y + 0.25, 0),
            FreeCAD.Vector(START+i, y - 0.25, 0)
        ]
        Draft.makeWire(points, closed=False, face=True, support=None)
        show()
        Draft.makeText(str(i+start_num), FreeCAD.Vector(START+i - 0.15, y - 0.7))
        show()

    # draw frame
    for i in range(0, int(END-START)+1, FRAME):
        points = [
            FreeCAD.Vector(START + i, y + 3, 0),
            FreeCAD.Vector(START + i, y, 0)
        ]
        Draft.makeWire(points, closed=False, face=True, support=None)
        show()


def draw_block(start, end, y, text):
    Gui.activateWorkbench("DrawingWorkbench")
    start = (START + start)
    end = (START + end)
    pl = FreeCAD.Placement()
    pl.Base = FreeCAD.Vector(start, y, 0.0)
    Draft.makeRectangle(length=end - start, height=2, placement=pl, face=False, support=None)
    show()
    Draft.makeText(text, FreeCAD.Vector(start + float(end - start)/2 - 0.15, y + 1))
    show()


# Create document
doc = App.newDocument(DOC_NAME)
App.ActiveDocument = doc

# Create drawing frame
doc.addObject('Drawing::FeaturePage', 'Page')
doc.Page.Template = TEMPLATE_PATH
Gui.activateWorkbench("DraftWorkbench")

draw_track(TRACK_Y, 0)
tracks = [TRACK_Y]
with open(INPUT_PATH) as f:
    tasks = [Task.decode(obj) for obj in json.load(f)]
    tasks = Task.merge(tasks)
    draw_table(tasks)
    for task in tasks:
        for period in task.periods:
            y = TRACK_Y
            if period["start"] != 0 and period["end"] != 0:
                y = TRACK_Y - TRACK_MARGIN*(period["start"] // ((END-START + 1)*1000))
                period["start"] = period["start"] % ((END-START + 1)*1000)
                period["end"] = period["end"] % ((END-START + 1)*1000)
                if y not in tracks:
                    draw_track(y, int(END-START) + len(tracks))
                    tracks.append(y)
            draw_block(float(period["start"])/1000, float(period["end"])/1000, y, task.id)

Gui.activateWorkbench("ArchWorkbench")

# python 2.7
# execfile('/home/egor/go/src/srv2/freecad.py')

# python3+
# exec(open("/home/egor/go/src/srv2/freecad.py").read())
