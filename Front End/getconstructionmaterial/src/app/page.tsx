import Image from 'next/image';
import { Analytics } from '@vercel/analytics/react';
import { SpeedInsights } from "@vercel/speed-insights/next";

import NavBar from './components/nav'
import ProductSubmissionComponent from './components/productSubmission'

export default function LandingPage() {
  return (

    <html>
      <head>
        <title>Get Construction Material</title>
        <link rel="canonical" href="https://www.getconstructionmaterial.com"/>
      </head>

      <body className = "">
          <NavBar />

          <main className="min-h-screen bg-brand flex flex-col items-center w-full pt-4 px-8">

            <div className="dark:text-black">
              <div className="font-mono">
                <h1 className="text-5xl sm:text-6xl pb-4">Looking for Construction Material?</h1>
                <p className="pb-4">Get Construction Material will provide you with an easy way to search what material is avaialble near you and get the best price possible.</p>
                
                <h2 className="pb-4 sm:pt-0 font-bold sm:text-xl">We will provide you with material information within 2 business days.</h2>
                <p className="sm:pt-0 pb-4">Tell us what you are looking for below and we will try our best to help with your construction needs.</p>
              </div>

              <ProductSubmissionComponent />

            </div>
            
          </main>
          <Analytics />
          <SpeedInsights />
        </body>


    </html>


    

    
    
  
    
  )
}
