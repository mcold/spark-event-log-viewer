package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"strconv"
	"strings"
)

type pagePlanType struct {
	root    *tview.TreeNode
	nNum    int
	mIdDesc map[int]string
	*tview.TreeView
	*tview.Flex
}

var pagePlan pagePlanType

func (pagePlan *pagePlanType) build() {

	pagePlan.mIdDesc = make(map[int]string)

	pagePlan.root = tview.NewTreeNode(".").
		SetColor(tcell.ColorRed)

	pagePlan.TreeView = tview.NewTreeView().SetRoot(pagePlan.root)
	pagePlan.TreeView.SetCurrentNode(pagePlan.root)
	pagePlan.TreeView.SetTitle("plan").
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetTitleAlign(tview.AlignLeft)

	setTreePlan()

	pagePlan.root.SetExpanded(true)

	pagePlan.Flex = tview.NewFlex().
		AddItem(pagePlan.TreeView, 0, 1, true)

	pagePlan.Flex.SetFocusFunc(func() {
		log.Println("SetFocusFunc")
		pagePlan.nNum = -1
		setTreePlan()
		app.SetFocus(pagePlan.TreeView)
	})
	// TODO
	//pagePlan.TreeView.SetSelectedFunc()

	application.pages.AddPage("plan", pagePlan.Flex, true, false)
}

func setTreePlan() {
	log.Println("setTreePlan")
	ev := pageMain.Events[pageMain.List.GetCurrentItem()]
	arrNums := getNodeNum(ev)

	pagePlan.nNum = 1
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
		log.Println(str)
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
