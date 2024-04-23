"use client";

import React, { useEffect, useRef, memo } from "react";
import * as d3 from "d3";

interface Node extends d3.SimulationNodeDatum {
  id: number;
  title: string;
}

const MyChart: React.FC = () => {
  const svgRef = useRef<SVGSVGElement>(null);

  useEffect(() => {
    console.log("MyChart");
    const margin = { top: 30, right: 80, bottom: 30, left: 30 };
    const width = 590;
    const height = 400;

    const svg = d3
      .select(svgRef.current)
      .attr("width", width + margin.left + margin.right)
      .attr("height", height + margin.top + margin.bottom)
      .append("g")
      .attr("transform", `translate(${margin.left},${margin.top})`);

    const dataset = {
      // make 100 nodes
      nodes: [
        { id: 1, title: "Node 1" },
        { id: 2, title: "Node 2" },
        { id: 3, title: "Node 3" },
        { id: 4, title: "Node 4" },
        { id: 5, title: "Node 5" },
        { id: 6, title: "Node 6" },
        { id: 7, title: "Node 7" },
      ],
      // links 1-91
      links: [
        { source: 1, target: 2 },
        { source: 2, target: 3 },
        { source: 2, target: 4 },
        { source: 3, target: 4 },
        { source: 4, target: 5 },
        { source: 5, target: 6 },
        { source: 6, target: 7 },
      ],
    };

    const simulation = d3
      .forceSimulation<Node, d3.SimulationLinkDatum<Node>>(dataset.nodes)
      .force(
        "link",
        d3
          .forceLink(dataset.links)
          .id((d: any) => d.id)
          .distance(100)
      )
      .force("charge", d3.forceManyBody())
      .force("center", d3.forceCenter(width / 2, height / 2));

    const arrowMarker = svg
      .append("defs")
      .append("marker")
      .attr("id", "arrow")
      .attr("viewBox", "0 -5 10 10")
      .attr("refX", 18)
      .attr("refY", 0)
      .attr("markerWidth", 20)
      .attr("markerHeight", 20)
      .attr("orient", "auto");

    const link = svg
      .append("g")
      .attr("class", "links")
      .selectAll("line")
      .data(dataset.links)
      .enter()
      .append("line")
      .style("stroke", "black")
      .attr("marker-end", "url(#arrow)");

    const node = svg
      .append("g")
      .attr("class", "nodes")
      .selectAll("circle")
      .data(dataset.nodes)
      .enter()
      .append("circle")
      .attr("r", 20)
      .call(
        d3
          .drag<SVGCircleElement, Node>()
          .on("start", (event: any, d: any) => dragstarted(event, d))
          .on("drag", (event: any, d: any) => dragged(event, d))
          .on("end", (event: any, d: any) => dragended(event, d))
      )
      .style("fill", "black");

    const text = svg
      .append("g")
      .attr("class", "text")
      .selectAll("text")
      .data(dataset.nodes)
      .enter()
      .append("text")
      .text((d: Node) => d.id)
      .style("fill", "white");

    const title = svg
      .append("g")
      .attr("class", "title")
      .selectAll("text")
      .data(dataset.nodes)
      .enter()
      .append("text")
      .text((d: Node) => d.title)
      .style("fill", "black")
      .style("font-size", "20px")
      .attr("class", "font-schoolbell");

    simulation.on("tick", ticked);

    function ticked() {
      link
        .attr("x1", (d: any) => d.source.x)
        .attr("y1", (d: any) => d.source.y)
        .attr("x2", (d: any) => d.target.x)
        .attr("y2", (d: any) => d.target.y);

      arrowMarker
        .append("path")
        .attr("d", "M0,-5L10,0L0,5")
        .attr("fill", "black");

      node.attr("cx", (d: any) => d.x).attr("cy", (d: any) => d.y);

      text.attr("x", (d: any) => d.x - 5).attr("y", (d: any) => d.y + 5);

      title.attr("x", (d: any) => d.x - 20).attr("y", (d: any) => d.y - 30);
    }

    function dragstarted(event: any, d: any) {
      if (!event.active) simulation.alphaTarget(0.3).restart();
      d.fx = d.x;
      d.fy = d.y;
    }

    function dragged(event: any, d: any) {
      d.fx = event.x;
      d.fy = event.y;
    }

    function dragended(event: any, d: any) {
      if (!event.active) simulation.alphaTarget(0);
      d.fx = null;
      d.fy = null;
    }

    return () => {
      simulation.stop();
      svg.selectAll("*").remove();
    };
  }, []);

  return <svg ref={svgRef} className='border-2 border-black' />;
};

export default memo(MyChart);
