import React from 'react';
import { motion } from 'framer-motion';
import Navbar from './components/Navbar';
import Hero from './components/Hero';
import Features from './components/Features';
import HowItWorks from './components/HowItWorks';
import Architecture from './components/Architecture';
import QuickStart from './components/QuickStart';
import CodeExamples from './components/CodeExamples';
import Comparison from './components/Comparison';
import Footer from './components/Footer';

function App() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-dark-900 via-dark-800 to-dark-900 text-white">
      <Navbar />
      <main>
        <Hero />
        <Features />
        <HowItWorks />
        <Architecture />
        <QuickStart />
        <CodeExamples />
        <Comparison />
      </main>
      <Footer />
    </div>
  );
}

export default App; 