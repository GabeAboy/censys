import React from 'react';

const Checkbox = ({ checked, onChange, label, colorClass = 'blue' }) => {
  const colorMap = {
    blue: 'text-blue-600 focus:ring-blue-500',
    red: 'text-red-600 focus:ring-red-500',
    yellow: 'text-yellow-600 focus:ring-yellow-500',
    green: 'text-green-600 focus:ring-green-500',
    gray: 'text-gray-600 focus:ring-gray-500'
  };

  const labelColorMap = {
    blue: 'text-gray-700',
    red: 'text-red-700',
    yellow: 'text-yellow-700',
    green: 'text-green-700',
    gray: 'text-gray-700'
  };

  return (
    <label className="flex items-center gap-2 cursor-pointer">
      <input
        type="checkbox"
        checked={checked}
        onChange={onChange}
        className={`w-4 h-4 border-gray-300 rounded focus:ring-2 ${colorMap[colorClass] || colorMap.blue}`}
      />
      <span className={`text-sm ${labelColorMap[colorClass] || labelColorMap.blue} font-medium`}>
        {label}
      </span>
    </label>
  );
};

export default Checkbox;
