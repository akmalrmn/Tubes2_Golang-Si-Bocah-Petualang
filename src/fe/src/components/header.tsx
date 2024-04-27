import React from "react";
import Image from "next/image";
import logo from "../../public/logo.png";
import Link from "next/link";

const route = {
  "/wikirace": "WikiRace",
  "/about-us": "About Us",
};

function Header() {
  return (
    <div className='top-0 h-20 w-full sticky -mt-20 flex items-center justify-between py-8 px-12 z-20'>
      <div>
        <Link href='/'>
          <Image
            src={logo}
            alt='logo'
            height={60}
            className='hover:scale-125 transition-all'
          />
        </Link>
      </div>

      <div className='flex space-x-10 font-schoolbell text-2xl'>
        {Object.entries(route).map(([path, name]) => (
          <Link
            key={path}
            href={path}
            className='hover:scale-125 transition-all'
          >
            {name}
          </Link>
        ))}
      </div>

      <div className='absolute -inset-1 bg-white blur-sm z-[-2]'></div>
      <div className='absolute -inset-[0.1rem] bg-white z-[-1]'></div>
    </div>
  );
}

export default Header;
