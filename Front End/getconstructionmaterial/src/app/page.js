import Image from 'next/image'

import NavBar from './components/nav'
import ProductSubmissionComponent from './components/productSubmission'

export default function LandingPage() {
  return (


    <body className = "">
      <NavBar />

      <main className="min-h-screen bg-brand flex flex-col items-center w-full pt-4 px-8">

        <div className="">
          <div className="font-mono">
            <h1 className="text-5xl sm:text-6xl pb-4">Looking for Construction Material?</h1>
            <p className="pb-4">Get Construction Material will provide you with an easy way to search what material is avaialble near you and get the best price possible.</p>
            
            <h2 className="pb-4 sm:pt-0 font-bold sm:text-xl">We will provide you with material information within 1 business day.</h2>
            <p className="sm:pt-0 pb-4">Tell us what you are looking for and provide your email.</p>
          </div>

          <ProductSubmissionComponent />

        </div>
      
      </main>
    </body>

    
    
  
    
  )
}
