import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import React from "react"
import ClearLogs from "./ClearLogs"
import { logLinesToString } from "./logs"
import LogStore, { LogStoreProvider } from "./LogStore"
import { appendLinesForManifestAndSpan } from "./testlogs"
import { ResourceName } from "./types"

describe("ClearLogs", () => {
  const createPopulatedLogStore = (): LogStore => {
    const logStore = new LogStore()
    appendLinesForManifestAndSpan(logStore, "", "", [
      "global 1\n",
      "global 2\n",
    ])
    appendLinesForManifestAndSpan(logStore, "vigoda", "build:m1:1", [
      "m1:1 build line 1\n",
    ])
    appendLinesForManifestAndSpan(logStore, "vigoda", "pod:m1-abc123", [
      "m1:1 runtime line 1\n",
    ])
    appendLinesForManifestAndSpan(logStore, "manifest2", "build:m2", [
      "m2 build line 1\n",
    ])
    appendLinesForManifestAndSpan(logStore, "vigoda", "build:m1:2", [
      "m1:2 build line 1\n",
      "m1:2 build line 2\n",
    ])
    appendLinesForManifestAndSpan(logStore, "manifest2", "pod:m2-def456", [
      "m2 runtime line 1\n",
    ])
    return logStore
  }

  it("clears all resources", () => {
    const logStore = createPopulatedLogStore()
    render(
      <LogStoreProvider value={logStore}>
        <ClearLogs resourceName={ResourceName.all} />
      </LogStoreProvider>
    )

    userEvent.click(screen.getByRole("button"))

    expect(logStore.spans).toEqual({})
    expect(logStore.allLog()).toHaveLength(0)
  })

  it("clears a specific resource", () => {
    const logStore = createPopulatedLogStore()
    render(
      <LogStoreProvider value={logStore}>
        <ClearLogs resourceName={"vigoda"} />
      </LogStoreProvider>
    )

    userEvent.click(screen.getByRole("button"))

    expect(Object.keys(logStore.spans).sort()).toEqual([
      "_",
      "build:m2",
      "pod:m2-def456",
    ])

    expect(logLinesToString(logStore.allLog(), false)).toEqual(
      "global 1\nglobal 2\nm2 build line 1\nm2 runtime line 1"
    )
  })
})
