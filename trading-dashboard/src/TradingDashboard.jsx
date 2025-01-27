import React, { useState, useEffect } from 'react';
import axios from 'axios';

const TradingDashboard = () => {
  const [stockData, setStockData] = useState(null);
  const [loading, setLoading] = useState(true);

  // Fetch stock data when the component mounts
  useEffect(() => {
    axios.get('https://api.example.com/stock-data')  // Replace with actual API URL
      .then(response => {
        setStockData(response.data);
        setLoading(false);
      })
      .catch(error => {
        console.error('Error fetching stock data', error);
        setLoading(false);
      });
  }, []);

  return (
    <div>
      <h1>Trading Dashboard</h1>
      {loading ? (
        <p>Loading data...</p>
      ) : (
        <div>
          {stockData ? (
            <ul>
              {stockData.map(stock => (
                <li key={stock.symbol}>
                  <strong>{stock.name}:</strong> {stock.price}
                </li>
              ))}
            </ul>
          ) : (
            <p>No stock data available</p>
          )}
        </div>
      )}
    </div>
  );
};

export default TradingDashboard;
