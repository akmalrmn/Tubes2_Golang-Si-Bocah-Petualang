import { ParticleBackground } from "@/components/particle-background";
import Image from "next/image";
import hero from "@/assets/images/hero.png";

export default function Home() {
  return (
    <main className='flex min-h-screen flex-col items-center justify-between p-24'>
      <div className='h-full flex flex-col items-center -space-y-10'>
        <div className='py-1 px-12 w-fit relative'>
          <div className='absolute -inset-1 bg-white blur-lg z-[-1] rounded-full'></div>
          <h1 className='font-schoolbell text-[100px] bg-white'>GoLang</h1>
        </div>

        <div className='py-8 px-12 w-fit relative'>
          <div className='absolute -inset-1 bg-white blur-lg z-[-1] rounded-full'></div>
          <h1 className='font-schoolbell text-[100px] bg-white'>
            Si Bocah Petualang
          </h1>
        </div>

        <div className='py-8 px-16 mt-2 w-fit relative'>
          <div className='absolute -inset-1 bg-white blur-lg z-[-1] rounded-full'></div>
          <h1 className='font-schoolbell text-center text-xl w-[100vh] leading-loose bg-white'>
            Welcome to WikiRace Solver! Utilizing Breadth-First Search (BFS) and
            Iterative Deepening Search (IDS) algorithms, our web app will find
            the shortest path between any two Wikipedia pages. Developed as part
            of the Algorithm Strategy course's second major assignment, it
            demonstrates our prowess in algorithm implementation and
            optimization.
          </h1>
        </div>
      </div>
      <ParticleBackground />
    </main>
  );
}
