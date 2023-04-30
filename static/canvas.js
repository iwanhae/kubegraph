import * as d3 from "https://cdn.jsdelivr.net/npm/d3@7/+esm";

function drag(simulation) {
    function dragstarted(event) {
        if (!event.active) simulation.alphaTarget(0.3).restart();
        event.subject.fx = event.subject.x;
        event.subject.fy = event.subject.y;
    }

    function dragged(event) {
        event.subject.fx = event.x;
        event.subject.fy = event.y;
    }

    function dragended(event) {
        if (!event.active) simulation.alphaTarget(0);
        event.subject.fx = null;
        event.subject.fy = null;
    }

    return d3.drag()
        .on("start", dragstarted)
        .on("drag", dragged)
        .on("end", dragended);
}

const mindistance = 30
const inrange = ({ x: sx, y: sy }, { x: tx, y: ty }) => Math.hypot(sx - tx, sy - ty) <= mindistance

export function graph(width = 600, height = 600) {
    const nodes = [];
    const links = [];
    let mouse = null;

    const svg = d3.create("svg")
        .property("value", { nodes: [], links: [] })
        .attr("viewBox", [-width / 2, -height / 2, width, height])
        .attr("cursor", "crosshair")
        .on("mouseleave", mouseleft)
        .on("mousemove", mousemoved)
        .on("click", clicked);

    const simulation = d3.forceSimulation(nodes)
        .force("charge", d3.forceManyBody().strength(-60))
        .force("link", d3.forceLink(links))
        .force("x", d3.forceX())
        .force("y", d3.forceY())
        .on("tick", ticked);

    const dragger = drag(simulation)
        .on("start.mouse", mouseleft)
        .on("end.mouse", mousemoved);

    let link = svg.append("g")
        .attr("stroke", "#999")
        .selectAll("line");

    let mouselink = svg.append("g")
        .attr("stroke", "red")
        .selectAll("line");

    let node = svg.append("g")
        .selectAll("circle");

    const cursor = svg.append("circle")
        .attr("display", "none")
        .attr("fill", "none")
        .attr("stroke", "red")
        .attr("r", mindistance - 5);

    function ticked() {
        node.attr("cx", d => d.x)
            .attr("cy", d => d.y)
            .attr("stroke", "none")
            .attr("fill", d => d.content.color ?? "black")

        link.attr("x1", d => d.source.x)
            .attr("y1", d => d.source.y)
            .attr("x2", d => d.target.x)
            .attr("y2", d => d.target.y);

        mouselink = mouselink
            .data(mouse ? nodes.filter(node => inrange(mouse, node)) : [])
            .join("line")
            .attr("x1", mouse && mouse.x)
            .attr("y1", mouse && mouse.y)
            .attr("x2", d => d.x)
            .attr("y2", d => d.y);

        cursor
            .attr("display", mouse ? null : "none")
            .attr("cx", mouse && mouse.x)
            .attr("cy", mouse && mouse.y);
    }

    function mouseleft() {
        mouse = null;
    }

    function mousemoved(event) {
        const [x, y] = d3.pointer(event);
        mouse = { x, y };
        document.getElementById("live").innerText = simulation.find(mouse.x, mouse.y).id
        simulation.alpha(0.1).restart();
    }

    function clicked(event) {
        mousemoved.call(this, event);
        simulation.alpha(0.5).restart();
        const tgt = simulation.find(mouse.x, mouse.y)
        console.log(tgt)
        document.getElementById("info").innerText = tgt.id
    }

    function remove(id) {
        let idx = nodes.findIndex((v) => v.id === id)
        if (0 <= idx) nodes.splice(idx, 1)

        while (true) {
            const idx = links.findIndex((v) => v.source.id === id || v.target.id === id)
            if (idx < 0) break
            links.splice(idx, 1)
        }

        simulation.nodes(nodes);
        simulation.force("link").links(links);
        simulation.alpha(0.01).restart();
    }

    function spawn(source) {
        const oldIdx = nodes.findIndex(v => v.id === source.id)
        if (0 <= oldIdx) {
            const old = nodes[oldIdx]
            source.x = old.x
            source.y = old.y
            if (0 <= oldIdx) nodes.splice(oldIdx, 1)
        }
        nodes.push(source);

        links.splice(0, links.length)
        nodes.forEach((source) => {
            if (!source.edges) return
            source.edges.forEach((id) => {
                const target = nodes.find(v => v.id === id)
                if (target !== undefined) links.push({ source, target })
            })
        })

        link = link
            .data(links)
            .join("line");

        node = node
            .data(nodes)
            .join(
                enter => enter.append("circle").attr("r", 0)
                    .call(enter => enter.transition().attr("r", 5))
                    .call(dragger),
                update => update,
                exit => exit.remove()
            );

        simulation.nodes(nodes);
        simulation.force("link").links(links);
        simulation.alpha(0.25).restart();
    }

    return Object.assign(svg.node(),
        { spawn, remove });
}