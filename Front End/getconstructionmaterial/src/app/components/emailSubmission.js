import React, {useState} from 'react';

function EmailSubmission() {
  const [formData, setFormData] = useState({
    email: ''
  });

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch('http://localhost:443/emailForm', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
      });
      if (response.ok) {
        console.log('Form submitted successfully');
      } else {
        console.error('Form submission failed');
      }
    } catch (error) {
      console.error('Error submitting form', error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col sm:flex-row sm:w-[600px]">
        <input type="email" name="email" value={formData.email} onChange={handleChange} placeholder="E-mail Address" className="sm:flex items-stretch flex-grow focus:outline-none block rounded-lg sm:rounded-none sm:rounded-l-lg pl-4 py-2"></input>
        
        <button type="submit" className="sm:mt-0 sm:w-auto sm:-ml-2 py-2 px-2 rounded-lg font-medium text-white focus:outline-none bg-logo-blue">
        Stay in the Loop
        </button>
    </form>
  );
}

export default EmailSubmission;
