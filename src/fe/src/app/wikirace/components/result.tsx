"use client";

import { useEffect, useState, useRef } from "react";
import { useSearchParams } from "next/navigation";
import Loading from "../loading";
import MyChart from "./graph";

interface ResultProps {
  start: string | null;
  target: string | null;
}

function fetchPath(start: string | null, target: string | null) {
  if (!start || !target) {
    return Promise.resolve({
      start: null,
      target: null,
    });
  }
  //artificially slow down the fetch
  return new Promise<ResultProps>((resolve) => {
    setTimeout(() => {
      resolve({
        start: start,
        target: target,
      });
    }, 3000);
  });
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
  const [data, setData] = useState<string>("");

  const fetchData = async () => {
    setIsLoading(true);
    const result = await fetchPath(start, target);
    setData(result.target?.toString() || "");
    setIsLoading(false);
  };

  useEffect(() => {
    const elem = document.getElementById("result");

    if (elem) {
      elem.scrollIntoView({ behavior: "smooth" });
    }
    fetchData();
  }, [start, target, algorithm]);

  if (!start || !target) {
    return null;
  }

  if (isLoading) {
    return <Loading {...props} />;
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
            <span>2</span>
          </div>

          <div className='flex flex-col gap-2 w-fit'>
            <span>Total article visited:</span>
            <span>2</span>
          </div>

          <div className='flex flex-col gap-2 w-fit'>
            <span>Time taken:</span>
            <span>2</span>
          </div>
        </div>

        <div className='bg-white'>
          <MyChart />
        </div>
      </div>
    </div>
  );
}

export default Result;
