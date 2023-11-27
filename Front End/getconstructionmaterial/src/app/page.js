import Image from 'next/image'

import NavBar from './components/nav'

export default function LandingPage() {
  return (


    <body className = "w-full min-h-screen">
      <NavBar />
      <main className="flex flex-col items-center w-full pt-4 px-8">

        <div className="">
          <div className="font-mono">
            <h1 className="text-5xl sm:text-6xl">Looking for Construction Material?</h1>
            <p className="">Get Construction Material will provide you with an easy way to search what material is avaialble near you and get the best price possible.</p>
            <p>Sign up to be alerted when you can start searching!</p>
          </div>

          <form method="post" className="flex flex-col sm:flex-row sm:justify-between sm:items-center w-full sm:w-4/5 md:w-120 rounded-lg">
              <input type="email" placeholder="E-mail Address" className="flex-grow rounded placeholder-gray-500"></input>
              <button class="border-4 border-black hover:bg-brand-bg text-black rounded">
                Stay in the Loop
              </button>
          </form>
        </div>
      
        

      
      </main>
    </body>

    
    
  
    
  )
}
