"use client";
import React from "react";
import LoadingBar from "./components/loading-bar";
import RandomText from "./components/random-text";

function Loading(
  props: React.DetailedHTMLProps<
    React.HTMLAttributes<HTMLDivElement>,
    HTMLDivElement
  >
) {
  return (
    <div
      {...props}
      className='min-h-screen min-w-[100vw] flex flex-col items-center justify-center bg-white z-10 relative gap-8'
    >
      <LoadingBar />
      <RandomText />
      <div className='w-full h-full absolute -inset-2 bg-white blur-md'></div>
    </div>
  );
}

export default Loading;
