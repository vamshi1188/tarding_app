import React, { useState, useEffect } from 'react';
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

const TradingDashboard = () => {
  const [marketData, setMarketData] = useState([]);
  const [tradingCalls, setTradingCalls] = useState([
    { id: 1, symbol: "AAPL", call: "Buy", target: 150, stopLoss: 140 },
    { id: 2, symbol: "GOOGL", call: "Sell", target: 2700, stopLoss: 2750 },
  ]);

  useEffect(() => {
    // Example API call to fetch market data
    const fetchMarketData = async () => {
      try {
        const response = await fetch("https://api.example.com/market-data");
        const data = await response.json();
        setMarketData(data);
      } catch (error) {
        console.error("Error fetching market data:", error);
      }
    };

    fetchMarketData();
  }, []);

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <h1 className="text-3xl font-bold mb-6">Trading Dashboard</h1>

      <div className="mb-8">
        <h2 className="text-xl font-semibold mb-4">Market Data</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {marketData.length ? (
            marketData.map((stock) => (
              <Card key={stock.symbol} className="p-4">
                <CardContent>
                  <h3 className="text-lg font-bold">{stock.symbol}</h3>
                  <p>Price: ${stock.price}</p>
                  <p>Change: {stock.change}%</p>
                </CardContent>
              </Card>
            ))
          ) : (
            <p>Loading market data...</p>
          )}
        </div>
      </div>

      <div>
        <h2 className="text-xl font-semibold mb-4">Trading Calls</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {tradingCalls.map((call) => (
            <Card key={call.id} className="p-4">
              <CardContent>
                <h3 className="text-lg font-bold">{call.symbol}</h3>
                <p>Call: {call.call}</p>
                <p>Target: ${call.target}</p>
                <p>Stop Loss: ${call.stopLoss}</p>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>

      <Button className="mt-6" onClick={() => alert('Trading Call Added!')}>
        Add New Trading Call
      </Button>
    </div>
  );
};

export default TradingDashboard;
