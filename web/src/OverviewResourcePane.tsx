import React, { useEffect, useState } from "react"
import SplitPane from "react-split-pane"
import styled from "styled-components"
import { Alert, combinedAlerts } from "./alerts"
import { ApiButtonType, buttonsForComponent } from "./ApiButton"
import HeaderBar, { HeaderBarPage } from "./HeaderBar"
import { LogUpdateAction, LogUpdateEvent, useLogStore } from "./LogStore"
import OverviewResourceDetails from "./OverviewResourceDetails"
import OverviewResourceSidebar from "./OverviewResourceSidebar"
import "./Resizer.scss"
import { useResourceNav } from "./ResourceNav"
import { useSidebarContext } from "./SidebarContext"
import StarredResourceBar, {
  starredResourcePropsFromView,
} from "./StarredResourceBar"
import { Color, Width } from "./style-helpers"
import { ResourceName, UIResource } from "./types"

type OverviewResourcePaneProps = {
  view: Proto.webviewView
  isSocketConnected: boolean
}

let OverviewResourcePaneRoot = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100vh;
  background-color: ${Color.gray20};
  max-height: 100%;
`
let Main = styled.div`
  display: flex;
  width: 100%;
  // In Safari, flex-basis "auto" squishes OverviewTabBar + OverviewResourceBar
  flex: 1 1 100%;
  overflow: hidden;
  position: relative;

  .SplitPane {
    position: relative !important;
  }
  .Pane {
    display: flex;
  }
`

export default function OverviewResourcePane(props: OverviewResourcePaneProps) {
  let nav = useResourceNav()
  const logStore = useLogStore()
  let resources = props.view?.uiResources || []
  let name = nav.invalidResource || nav.selectedResource || ""
  let r: UIResource | undefined
  let all = name === "" || name === ResourceName.all
  let starred = name === ResourceName.starred
  if (!all) {
    r = resources.find((r) => r.metadata?.name === name)
  }
  let selectedTab = ""
  if (all) {
    selectedTab = ResourceName.all
  } else if (starred) {
    selectedTab = ResourceName.starred
  } else if (r?.metadata?.name) {
    selectedTab = r.metadata.name
  }

  const { isSidebarOpen, setSidebarOpen, setSidebarClosed } =
    useSidebarContext()

  const [paneSize, setPaneSize] = useState<number>(
    isSidebarOpen ? Width.sidebarDefault : Width.sidebarMinimum
  )

  // listen for changes from sidebar context in case it was toggled instead
  // being dragged past a breakpoint.
  useEffect(() => {
    setPaneSize(
      isSidebarOpen ? Width.sidebarDefault : Width.sidebarMinimum + 0.01
      // adds 0.01 so there's still a state diff when the user releases after dragging
    )
  }, [isSidebarOpen])

  const [truncateCount, setTruncateCount] = useState<number>(0)

  // add a listener to rebuild alerts whenever a truncation event occurs
  // truncateCount is a dummy state variable to trigger a re-render to
  // simplify logic vs reconciliation between logStore + props
  useEffect(() => {
    const rebuildAlertsOnLogClear = (e: LogUpdateEvent) => {
      if (e.action === LogUpdateAction.truncate) {
        setTruncateCount(truncateCount + 1)
      }
    }

    logStore.addUpdateListener(rebuildAlertsOnLogClear)
    return () => logStore.removeUpdateListener(rebuildAlertsOnLogClear)
  }, [truncateCount])

  let alerts: Alert[] = []
  if (r) {
    alerts = combinedAlerts(r, logStore)
  } else if (all) {
    resources.forEach((r) => alerts.push(...combinedAlerts(r, logStore)))
  }

  const buttons = buttonsForComponent(
    props.view.uiButtons,
    ApiButtonType.Resource,
    name
  )

  const handleSplitPaneResize = (newSize: number) => {
    if (newSize < Width.sidebarBreakpoint && isSidebarOpen) {
      setSidebarClosed()
    } else if (newSize >= Width.sidebarBreakpoint && !isSidebarOpen) {
      setSidebarOpen()
    }
  }

  return (
    <OverviewResourcePaneRoot>
      <HeaderBar
        view={props.view}
        currentPage={HeaderBarPage.Detail}
        isSocketConnected={props.isSocketConnected}
      />
      <StarredResourceBar
        {...starredResourcePropsFromView(props.view, selectedTab)}
      />
      <Main>
        <SplitPane
          split="vertical"
          size={paneSize}
          minSize={Width.sidebarMinimum}
          onChange={handleSplitPaneResize}
          onDragFinished={() =>
            setPaneSize(
              isSidebarOpen ? Width.sidebarDefault : Width.sidebarMinimum
            )
          }
        >
          <OverviewResourceSidebar {...props} name={name} />
          <OverviewResourceDetails
            resource={r}
            name={name}
            alerts={alerts}
            buttons={buttons}
          />
        </SplitPane>
      </Main>
    </OverviewResourcePaneRoot>
  )
}
