"use client";

import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import Loading from "../loading";

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

  if (true) {
    return <Loading {...props} />;
  }

  return (
    <div {...props} className='min-h-screen'>
      <h1>{data}</h1>
    </div>
  );
}

export default Result;
