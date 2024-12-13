'use client';

import { useState, useEffect, BaseSyntheticEvent } from 'react';
import Description from '@/app/_components/description';

export default function Goal() {
    const [message, setMessage] = useState('');
    const [previousMessages, setPreviousMessages] = useState([]);
    const [isDropdownVisible, setIsDropdownVisible] = useState(false);

    useEffect(() => {
        const getMessage = async () => {
            try {
                // const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/d2/terror/goal`);
                const response = await fetch(`/api/v1/d2/terror/goal`);
                if (!response.ok) {
                    throw new Error('Failed to fetch message');
                }
                const result = await response.json();
                setMessage(result[0].message);
                setPreviousMessages(result);
            } catch (error) {
                console.error('Error fetching message:', error);
            }
        };

        getMessage();
    }, []);

    const handleInputChange = (e: BaseSyntheticEvent) => {
        setMessage(e.target.value as string);
    };

    const handleSubmit = async (e: BaseSyntheticEvent) => {
        e.preventDefault();

        try {
            // const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/d2/terror/goal`, {
            const res = await fetch(`/api/v1/d2/terror/goal`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ message }),
            });
            alert('updated 😊');
        } catch (error) {
            alert(`failed to update 😔\nError: ${error}`);
        }
    };

    const handleOptionClick = (option) => {
        setMessage(option);
        setIsDropdownVisible(false);
    };

    return (
        <form className="w-full h-fit max-h-fit p-6 rounded-lg shadow-lg" onSubmit={handleSubmit}>
            <Description
                title=""
                subtitle={
                    <>
                        Move <strong>&#91;subs&#93;</strong> to where you want the number to be
                    </>
                }
            />
            <div className="flex items-center space-x-4 w-full mx-auto mb-4">
                <div className="relative w-full">
                    <input
                        value={message}
                        onChange={handleInputChange}
                        onFocus={() => setIsDropdownVisible(true)}
                        onBlur={() => setTimeout(() => setIsDropdownVisible(false), 200)}
                        placeholder="Update the !goal command here"
                        className="w-full px-4 py-2 text-gray-700 bg-gray-100 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-400"
                    />
                    {isDropdownVisible && (
                        <div className="absolute top-full left-0 z-10 mt-1 bg-white border border-gray-300 rounded shadow-lg max-h-50 overflow-y-auto w-full">
                            {previousMessages.map((msg, index) => (
                                <div
                                    key={index}
                                    onClick={() => handleOptionClick(msg.message)}
                                    className="px-4 py-2 cursor-pointer hover:bg-gray-200 text-black"
                                >
                                    {msg.message}
                                </div>
                            ))}
                        </div>
                    )}
                </div>
                <button
                    type="submit"
                    className="px-4 py-2 text-white bg-gradient-to-r from-primary to-secondary rounded hover:opacity-90 focus:outline-none focus:ring-2 focus:ring-purple-500"
                >
                    Submit
                </button>
            </div>
        </form>
    );
}
