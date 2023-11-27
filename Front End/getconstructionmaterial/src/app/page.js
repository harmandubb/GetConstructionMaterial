import Image from 'next/image'

import NavBar from './components/nav'

export default function LandingPage() {
  return (


    <body className = "">
      <NavBar />

      <main className="min-h-screen bg-brand flex flex-col items-center w-full pt-4 px-8">

        <div className="">
          <div className="font-mono">
            <h1 className="text-5xl sm:text-6xl pb-4">Looking for Construction Material?</h1>
            <p className="">Get Construction Material will provide you with an easy way to search what material is avaialble near you and get the best price possible.</p>
            <p className="pb-4 sm:pt-0 pt-4">Sign up to be alerted when you can start searching!</p>
          </div>

          

        </div>
      
      </main>
    </body>

    
    
  
    
  )
}
