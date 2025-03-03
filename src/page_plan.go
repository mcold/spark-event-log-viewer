package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
	"strings"
)

type pagePlanType struct {
	root      *tview.TreeNode
	nNum      int
	mIdDesc   map[int]string
	mIdSimple map[int]string
	*tview.TreeView
	simpleArea *tview.TextArea
	descArea   *tview.TextArea
	flDesc     *tview.Flex
	*tview.Flex
}

var pagePlan pagePlanType

func (pagePlan *pagePlanType) build() {

	pagePlan.mIdDesc = make(map[int]string)
	pagePlan.mIdSimple = make(map[int]string)

	pagePlan.root = tview.NewTreeNode(".").
		SetColor(tcell.ColorRed).
		SetSelectable(false)

	pagePlan.TreeView = tview.NewTreeView().SetRoot(pagePlan.root)
	pagePlan.TreeView.SetCurrentNode(pagePlan.root)
	pagePlan.TreeView.SetTitle("plan").
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetTitleAlign(tview.AlignLeft)

	setTreePlan()

	pagePlan.root.SetExpanded(true)

	pagePlan.simpleArea = tview.NewTextArea()

	pagePlan.simpleArea.SetBorder(true).SetBorderColor(tcell.ColorBlue)
	pagePlan.simpleArea.SetTitle("simple").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1)

	pagePlan.descArea = tview.NewTextArea()

	pagePlan.descArea.SetBorder(true).SetBorderColor(tcell.ColorBlue)
	pagePlan.descArea.SetTitle("description").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 1, 1)

	pagePlan.flDesc = tview.NewFlex().SetDirection(tview.FlexRow)
	pagePlan.flDesc.SetBorder(true).SetBorderColor(tcell.ColorBlue)

	pagePlan.flDesc.
		AddItem(pagePlan.simpleArea, 0, 1, true).
		AddItem(pagePlan.descArea, 0, 1, true)

	pagePlan.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(pagePlan.TreeView, 0, 5, true).
		AddItem(pagePlan.flDesc, 0, 3, true)

	pagePlan.Flex.SetFocusFunc(func() {
		pagePlan.nNum = -1
		setTreePlan()
		app.SetFocus(pagePlan.TreeView)
	})

	pagePlan.TreeView.SetSelectedFunc(
		func(node *tview.TreeNode) {
			pagePlan.descArea.SetText(pagePlan.mIdDesc[node.GetReference().(int)], true)
			pagePlan.simpleArea.SetText(pagePlan.mIdSimple[node.GetReference().(int)], true)
		})

	application.pages.AddPage("plan", pagePlan.Flex, true, false)
}

func setTreePlan() {
	log.Println("setTreePlan")
	log.Println("clear tree")
	pagePlan.root.ClearChildren()
	ev := pageMain.Events[pageMain.List.GetCurrentItem()]
	arrNums := getNodeNum(ev)

	pagePlan.nNum = 1
	pagePlan.mIdSimple[arrNums[pagePlan.nNum]] = ev.SparkPlanInfo.SimpleString
	node := tview.NewTreeNode(ev.SparkPlanInfo.NodeName + " (" + strconv.Itoa(arrNums[pagePlan.nNum]) + ")").
		SetSelectable(true).
		SetColor(tcell.ColorBlue).
		SetReference(arrNums[pagePlan.nNum]).
		SetExpanded(true)
	pagePlan.root.AddChild(node)
	setTreePlanChildren(node, ev.SparkPlanInfo, arrNums)
}

func setTreePlanChildren(rootNode *tview.TreeNode, plan SparkPlan, arrNums map[int]int) {
	var node *tview.TreeNode
	for _, child := range plan.Children {
		pagePlan.nNum++
		pagePlan.mIdSimple[arrNums[pagePlan.nNum]] = child.SimpleString
		node = tview.NewTreeNode(child.NodeName + " (" + strconv.Itoa(arrNums[pagePlan.nNum]) + ")").
			SetSelectable(true).
			SetColor(tcell.ColorBlue).
			SetReference(arrNums[pagePlan.nNum]).
			SetExpanded(true)
		rootNode.AddChild(node)
		setTreePlanChildren(node, child, arrNums)
	}
}

func getNodeNum(event Event) map[int]int {
	arrStr := strings.Split(event.PhysicalPlanDescription, "\n")
	arrRes := make(map[int]int)
	pagePlan.nNum = -1
	for _, str := range arrStr {
		pagePlan.nNum++
		if pagePlan.nNum == 0 {
			continue
		}
		if strings.TrimSpace(str) == "" {
			break
		}

		parts := strings.FieldsFunc(str, func(r rune) bool {
			return r == '(' || r == ')' || r == ' '
		})
		match := parts[len(parts)-1]
		if len(match) == 0 {
			break
		} else {
			id, err := strconv.Atoi(match)

			check(err)
			arrRes[pagePlan.nNum] = id
		}

	}

	// get description for each operation
	n := 1
	var arrTemp []string
	bFound := false
	for i, row := range arrStr[pagePlan.nNum+1:] {
		if strings.HasPrefix(row, "("+strconv.Itoa(n)+")") {
			bFound = true
			continue
		}
		// found next one
		if strings.HasPrefix(row, "("+strconv.Itoa(n+1)+")") {
			if len(arrTemp) > 0 {
				pagePlan.mIdDesc[n] = strings.Join(arrTemp, "\n")
				arrTemp = make([]string, 0)
			}
			n++
		} else {
			if bFound && strings.TrimSpace(row) != "" {
				arrTemp = append(arrTemp, row)
			}
		}

		if i == len(arrStr[pagePlan.nNum+1:])-1 && bFound && len(arrTemp) > 0 {
			pagePlan.mIdDesc[n] = strings.Join(arrTemp, "\n")
		}

	}
	return arrRes
}
