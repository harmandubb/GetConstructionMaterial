import Image from 'next/image'

import NavBar from '../../components/nav'

export default function LandingPage() {
  return (


    <body className = "bg-brand-bg">
      <NavBar />
      <main className="flex min-h-screen flex-col items-center justify-between p-24">
      
      <div className="font-mono px-4">
        <h1 className="text-6xl">Looking for Construction Material?</h1>
        <p className="py-4">Get Construction Material will provide you with an easy to search what material is avaialble near you and get the best price possible.</p>
        <p>Sign up to be alerted when you can start searching!</p>
      </div>

      <form method="post" class="bg-logo-blue p-4 flex justify-between items-center rounded-lg">
          <input type="email" placeholder="E-mail Address" class="flex-grow mr-4 p-2 rounded placeholder-gray-500"></input>
          <button class="bg-brand-red-dark hover:bg-logo-red-light text-white font-bold py-2 px-4 rounded">
            Stay in the Loop
          </button>
      </form>

      
    </main>
    </body>

    
    
  
    
  )
}
