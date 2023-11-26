import React from 'react';
import Image from 'next/image'

function NavBar() {
  return (
    <nav className="flex items-center">
      <Image 
        src = "/Logo/logo.png"
        width = {100}
        height= {100}
        alt = "Logo"
      />
      <div className="flex flex-col items-start">
          <p className="font-mono text-2xl font-black">Get Construction Material</p>
          <p className="font-mono text-sm font-black">by Docstruction</p>
      </div>

      
      
    </nav>
  );
}

export default NavBar;
