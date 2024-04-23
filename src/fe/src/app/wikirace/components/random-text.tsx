"use client";

import React, { useEffect, useState } from "react";

const randomTextLoading = [
  "Do you know that the shortest path between two points is a straight line?",
  "To get from Jokowi to Prabowo, you need to go through a lot of people",
  "The shortest path between two pages is not always the most obvious one",
];

function RandomText() {
  const [randomText, setRandomText] = useState<string>("Loading...");

  useEffect(() => {
    const interval = setInterval(() => {
      setRandomText(
        randomTextLoading[Math.floor(Math.random() * randomTextLoading.length)]
      );
    }, 5000);
    return () => clearInterval(interval);
  }, []);
  return (
    <p className='z-10 font-schoolbell text-xl transition-all'>{randomText}</p>
  );
}

export default RandomText;
