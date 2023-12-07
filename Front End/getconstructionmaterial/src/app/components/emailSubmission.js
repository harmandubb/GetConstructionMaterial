'use client'

import React, {useState} from 'react';




function EmailSubmission() {
  process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';

  const [formData, setFormData] = useState('');

  const [messageStatus, setMessageStatus] = useState({
    errorMessage: '',
    success: false  
  });


  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    console.log('Submitting form data:', formData); // Log form data
    if (formData.email != "") {
    try {
      const response = await fetch('https://api.getconstructionmaterial.com/emailFrom', { //Change for production
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
      });

      console.log('Response:', response); // Log response
      if (response.ok) {
        const responseBody = await response.json();
        // console.log('Response Body:', responseBody);
        setMessageStatus({message: "", success: true})
        // console.log('Form submitted successfully');
      } else {
        console.log("IN error visual handle")
        setMessageStatus({message:data.error, success:false})
        console.error('Form submission failed');
      }
    } catch (error) {
      setMessageStatus({message: "Network Error: Please try again shortly", success: false})
      console.log("Printing Errors:", error.message);

    }
  }
};


  return (
    <div className="flex flex-col sm:w-[600px]">
    <form onSubmit={handleSubmit} className="flex flex-col sm:flex-row">
        <input type="email" name="email" value={formData.email} onChange={handleChange} placeholder="E-mail Address" className="sm:flex items-stretch flex-grow focus:outline-none block rounded-lg sm:rounded-none sm:rounded-l-lg pl-4 py-2"></input>
        
        <button type="submit" className="sm:mt-0 sm:w-auto sm:-ml-2 py-2 px-2 rounded-lg font-medium text-white focus:outline-none bg-logo-blue">
        Stay in the Loop
        </button>
        
    </form>
      {messageStatus.message && <div className="border-2 rounded border-red-700 bg-red-300 py-2 px-2 mt-1">
        {messageStatus.message}</div>}
      {messageStatus.success && <div className="border-2 rounded border-green-700 bg-green-300 py-2 px-2 mt-1">
        Congrats! You have signed up. Stay Tuned!</div>}
    </div>
   
  );
}

export default EmailSubmission;
