"use client";

import { useEffect, useMemo, useState } from "react";
import Particles, { initParticlesEngine } from "@tsparticles/react";
import { type Container, type ISourceOptions } from "@tsparticles/engine";
import { loadFull } from "tsparticles"; // if you are going to use `loadFull`, install the "tsparticles" package too.

export const ParticleBackground = () => {
  const [init, setInit] = useState(false);

  useEffect(() => {
    initParticlesEngine(async (engine) => {
      await loadFull(engine);
    }).then(() => {
      setInit(true);
    });
  }, []);

  const particlesLoaded = async (container?: Container): Promise<void> => {
    console.log(container);
  };

  const configs: ISourceOptions = useMemo(
    () => ({
      fpsLimit: 120,
      fullScreen: {
        enable: true,
        zIndex: -1,
      },
      interactivity: {
        detect_on: "canvas",
        resize: {
          enable: true,
          mode: "bounce",
          density: 1,
          area: 800,
        },
        modes: {
          bubble: {
            distance: 400,
            duration: 2,
            opacity: 0.8,
            size: 40,
            speed: 3,
          },
          grab: { distance: 400, links: { opacity: 1 } },
          push: { quantity: 4 },
          remove: { quantity: 2 },
          repulse: { distance: 200, duration: 0.4 },
        },
      },
      particles: {
        color: { value: "#000000" },
        move: {
          attract: {
            enable: false,
            rotate: {
              x: 600,
              y: 1200,
            },
          },
          direction: "none",
          enable: true,
          outModes: "out",
          random: false,
          speed: 2,
          straight: false,
        },
        rotate: {
          animation: {
            enable: true,
            speed: 10,
            sync: false,
          },
        },
        number: { value: 120 },
        opacity: {
          animation: {
            enable: true,
            startValue: "min",
            speed: 0.1,
            sync: false,
            mode: "decrease",
          },
        },
        shape: {
          type: "character",
          options: {
            character: [
              {
                value:
                  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789".split(
                    ""
                  ),
                font: "Verdana",
                style: "",
                weight: "400",
                size: 60,
                fill: true,
              },
            ],
          },
        },
        size: {
          value: { min: 10, max: 30 },
          animation: {
            enable: true,
            speed: 2,
          },
        },
      },
      detectRetina: true,
    }),
    []
  );

  if (init) {
    return (
      <Particles
        id='tsparticles'
        particlesLoaded={particlesLoaded}
        options={configs}
        // add fade in animation
        className='animate-fade'
      />
    );
  }

  return <></>;
};
