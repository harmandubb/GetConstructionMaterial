import React, {useState} from 'react';

function emailSubmission() {
  return (
    <form method="post" className="flex flex-col sm:flex-row sm:w-[600px]">
        <input type="email" placeholder="E-mail Address" className="sm:flex items-stretch flex-grow focus:outline-none block rounded-lg sm:rounded-none sm:rounded-l-lg pl-4 py-2"></input>
        
        <button className="sm:mt-0 sm:w-auto sm:-ml-2 py-2 px-2 rounded-lg font-medium text-white focus:outline-none bg-logo-blue">
        Stay in the Loop
        </button>
    </form>
  );
}

export default NavBar;