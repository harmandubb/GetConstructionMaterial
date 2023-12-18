'use client'

import React, { useState, ChangeEvent, FormEvent } from 'react';

interface FormData {
    material: string;
    email: string;
    loc: string; 
}

interface MessageStatus {
  message?: string; 
  success?: boolean;
}

interface ErrorResponse {
  error: string;
}

const ProductSubmissionComponent: React.FC = () => {
  const [formData, setFormData] = useState<FormData>({email: '', material: '', loc: ''});
  const [messageStatus, setMessageStatus] = useState<MessageStatus>({});

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [event.target.name]: event.target.value });
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0'; //Should remove if not testing locally. 
    event.preventDefault();

    if (formData.email != "") {
      try {
        const response = await fetch('https://api.getconstructionmaterial.com/materialForm', { //Change for production
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(formData)
        });
  
        console.log('Response:', response); // Log response
        if (response.ok) {
          const responseBody = await response.json();
          console.log('Response Body:', responseBody);
          setMessageStatus({message: "", success: true})
          // console.log('Form submitted successfully');
        } else {
          const errorResponse: ErrorResponse = await response.json(); 
          setMessageStatus({message:errorResponse.error, success:false})
          console.error('Form submission failed');
        }
      } catch (error :any) {
        setMessageStatus({message: "Network Error: Please try again shortly", success: false})
        console.log("Printing Errors:", error.message);
  
      }
    }
  };

  return (
    <div className="flex flex-col lg:w-[1000px]">
      <form onSubmit={handleSubmit} className="flex flex-col lg:flex-row">
      <input type="search" name="material" value={formData.material} onChange={handleChange} placeholder="Material/Product wanted" className="sm:flex items-stretch flex-grow lg:border-r border-b-2 lg:border-b-0 focus:outline-none block rounded-lg lg:rounded-none lg:rounded-l-lg pl-4 py-2"></input>
      <input type="email" name="email" value={formData.email} onChange={handleChange} placeholder="E-mail Address" className="sm:flex items-stretch flex-grow sm:border-r border-b-2 lg:border-b-0 focus:outline-none block rounded-lg lg:rounded-none pl-4 py-2"></input>
      <input type="search" name="loc" value={formData.loc} onChange={handleChange} placeholder="City and Province" className="sm:flex items-stretch flex-grow focus:outline-none rounded-lg lg:rounded-none block pl-4 py-2"></input>  

        <button type="submit" className="sm:mt-0 sm:w-auto sm:-ml-2 py-2 px-2 rounded-lg font-medium text-white focus:outline-none bg-logo-blue">
            Find Material
        </button>
        
      </form>

      {messageStatus.message && <div className="border-2 rounded border-red-700 bg-red-300 py-2 px-2 mt-1">
        {messageStatus.message}</div>}
      {messageStatus.success && <div className="border-2 rounded border-green-700 bg-green-300 py-2 px-2 mt-1">
        Congrats! You have signed up. Stay Tuned!</div>}
    </div>
  );

}

export default ProductSubmissionComponent;
