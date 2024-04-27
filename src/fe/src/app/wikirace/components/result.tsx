"use client";

import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import Loading from "../loading";
import MyChart from "./graph";
import { SimulationLinkDatum, SimulationNodeDatum } from "d3-force";
import axios from "axios";

export interface Node extends SimulationNodeDatum {
  id: number;
  title: string;
}

export interface GraphData {
  nodes: Node[];
  links: SimulationLinkDatum<Node>[];
}

export interface resultData {
  time: number;
  totalArticleChecked: number;
  totalArticleVisited: number;
  graph: GraphData;
}

const dummyGraph: GraphData = {
  nodes: [
    { id: 1, title: "Node 1" },
    { id: 2, title: "Node 2" },
    { id: 3, title: "Node 3" },
    { id: 4, title: "Node 4" },
    { id: 5, title: "Node 5" },
  ],
  links: [
    { source: 1, target: 2 },
    { source: 1, target: 3 },
    { source: 1, target: 4 },
    { source: 4, target: 5 },
    { source: 2, target: 5 },
    { source: 3, target: 5 },
  ],
};

async function fetchPath(
  start: string | null,
  target: string | null,
  algorithm: string | null
): Promise<resultData> {
  try {
    const data = await axios
      .get(process.env.BACKEND_API_URL + "/" + algorithm, {
        params: {
          start,
          target,
        },
      })
      .then((res) => res.data);

    return data as resultData;
  } catch (e) {
    console.error(e);
    return {
      time: -1,
      totalArticleChecked: 0,
      totalArticleVisited: 0,
      graph: {
        nodes: [],
        links: [],
      },
    };
  }
}

function Result(
  props: React.DetailedHTMLProps<
    React.HTMLAttributes<HTMLDivElement>,
    HTMLDivElement
  >
) {
  const params = useSearchParams();

  const start = params.get("start");
  const target = params.get("target");
  const algorithm = params.get("algorithm");

  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [data, setData] = useState<resultData>({
    time: -1,
    totalArticleChecked: 0,
    totalArticleVisited: 0,
    graph: {
      nodes: [],
      links: [],
    },
  });

  const fetchData = async () => {
    setIsLoading(true);
    const result = await fetchPath(start, target, algorithm);
    setData(result);
    setIsLoading(false);
  };

  useEffect(() => {
    const elem = document.getElementById("result");

    if (elem) {
      elem.scrollIntoView({ behavior: "smooth" });
    }

    if (start && target && algorithm) {
      fetchData();
    }
  }, [start, target, algorithm]);

  if (!start || !target) {
    return null;
  }

  if (isLoading) {
    return <Loading {...props} />;
  }

  if (data.time === -1) {
    return null;
  }

  return (
    <div
      {...props}
      className='min-h-screen  w-[100vw] flex flex-col items-center z-20 p-20 gap-[5rem] font-schoolbell'
    >
      <div className='w-fit h-fit bg-white relative p-4'>
        <div className='bg-white -inset-10 blur-md absolute z-[-1]'></div>
        <h1 className='text-6xl z-30'>Result</h1>
      </div>
      <div className='flex justify-between items-center h-full bg-white relative -inset-4'>
        <div className='bg-white -inset-10 blur-md absolute z-[-1]'></div>
        <div className='flex flex-col justify-center flex-1 text-4xl gap-8 text-left mr-28 bg-white w-fit h-full'>
          <div className='flex flex-col gap-2 w-fit'>
            <span>Total article checked:</span>
            <span>{data.totalArticleChecked}</span>
          </div>

          <div className='flex flex-col gap-2 w-fit'>
            <span>Total article visited:</span>
            <span>{data.totalArticleVisited}</span>
          </div>

          <div className='flex flex-col gap-2 w-fit'>
            <span>Time taken:</span>
            <span>{data.time}</span>
          </div>
        </div>

        <div className='bg-white'>
          <MyChart dataset={data.graph} />
        </div>
      </div>
    </div>
  );
}

export default Result;
