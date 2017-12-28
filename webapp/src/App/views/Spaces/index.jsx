import React, { Component } from 'react';
import { Responsive as ResponsiveReactGridLayout } from "react-grid-layout";
import PkgSelector from '../../PkgSelector';
import Pkg from '../../Pkg';
import { AutoSizer } from 'react-virtualized';

function getStateFromLS(key) {
    if (global.localStorage) {
        return JSON.parse(global.localStorage.getItem("state")) || null;
    }
    return null;
}

function saveStateToLS(state) {
    if (global.localStorage) {
        global.localStorage.setItem(
            "state",
            JSON.stringify(state)
        );
    }
}

export default class Spaces extends Component {
    static defaultProps = {
        className: "layout",
        cols: { lg: 12, md: 10, sm: 6, xs: 4, xxs: 2 },
        rowHeight: 100
    };

    constructor(props) {
        super(props);

        this.state =
            getStateFromLS()
            ||
            {
                layouts: {},
                ops: [],
                newCounter: 0
            };
    }

    createElement(el) {
        return (
            <div key={el.i} data-grid={el}>
                <AutoSizer>
                    {({ height, width }) => (
                        <div style={{ height, width, border: 'dashed 3px #ececec', overflow: 'auto' }}>
                            <Pkg value={el.pkg} pkgRef={el.pkgRef} />
                        </div>
                    )}
                </AutoSizer>
                <span
                    className="remove"
                    style={{
                        position: "absolute",
                        right: "2px",
                        top: 0,
                        cursor: "pointer"
                    }}
                    onClick={this.onRemoveItem.bind(this, el.i)}
                >
                    x
                </span>
            </div>
        );
    }

    handlePkgAdded({ pkg, pkgRef }) {
        const op = {
            pkgRef,
            pkg,
            i: "n" + this.state.newCounter,
            x: (this.state.ops.length * 2) % (12),
            y: 10000000000000, // puts it at the bottom
            w: 2,
            h: 2
        }

        this.setState(
            prevState => ({
                ops: [...prevState.ops, op],
                newCounter: prevState.newCounter + 1
            })
        );
    }

    componentDidUpdate() {
        saveStateToLS(this.state);
    }

    onLayoutChange(layout, layouts) {
        this.setState({ layouts });
    }

    onRemoveItem(i) {
        this.setState(
            prevState => ({ ops: prevState.ops.filter(item => item.i !== i) })
        );
    }

    render() {
        return (
            <div>
                <PkgSelector
                    onSelect={selection => this.handlePkgAdded(selection)}
                />
                <AutoSizer>
                    {({ width }) => (
                        <ResponsiveReactGridLayout
                            width={width}
                            layouts={this.state.layouts}
                            onLayoutChange={(layout, layouts) => this.onLayoutChange(layout, layouts)}
                        >
                            {this.state.ops.map(el => (this.createElement(el)))}
                        </ResponsiveReactGridLayout>
                    )}
                </AutoSizer>
            </div>
        )
    }
}