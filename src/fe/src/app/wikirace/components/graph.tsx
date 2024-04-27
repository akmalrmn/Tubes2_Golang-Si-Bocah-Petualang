"use client";

import React, { useEffect, useRef, memo } from "react";
import * as d3 from "d3";
import { Node, GraphData } from "./result";

const MyChart = ({ dataset }: { dataset: GraphData }) => {
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

    const simulation = d3
      .forceSimulation<Node, d3.SimulationLinkDatum<Node>>(dataset.nodes)
      .force(
        "link",
        d3
          .forceLink(dataset.links)
          .id((d: any) => d.id)
          .distance(150)
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
      .attr("r", (d: Node) => {
        if (d.id === 1) return 25;
        if (d.id === dataset.nodes.length) return 25;
        return 20;
      })
      .style("fill", (d: Node) => {
        if (d.id === 1) return "#f05c5c";
        if (d.id === dataset.nodes.length) return "#6fa8dc";
        return "black";
      })
      .call(
        d3
          .drag<SVGCircleElement, Node>()
          .on("start", (event: any, d: any) => dragstarted(event, d))
          .on("drag", (event: any, d: any) => dragged(event, d))
          .on("end", (event: any, d: any) => dragended(event, d))
      )
      .style("cursor", "pointer");

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

      title.attr("x", (d: any) => d.x - 20).attr("y", (d: any) => d.y - 30);

      text.attr("x", (d: any) => d.x - 5).attr("y", (d: any) => d.y + 5);
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
