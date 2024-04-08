import React from "react";
import Image from "next/image";
import logo from "../../public/logo.png";

function Header() {
  return (
    <div className='top-0 h-20 w-full absolute flex items-center justify-between py-8 px-12'>
      <div>
        <Image src={logo} alt='logo' height={60} />
      </div>

      <div className='flex space-x-10 font-schoolbell text-2xl'>
        <div>Wikirace</div>

        <div>About Us</div>
      </div>

      <div className='absolute -inset-1 bg-white blur-sm z-[-2]'></div>
      <div className='absolute -inset-[0.1rem] bg-white z-[-1]'></div>
    </div>
  );
}

export default Header;
