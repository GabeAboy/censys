import React from 'react';

const AssetCard = ({ 
  asset, 
  editingTagForAsset, 
  newTagValue, 
  onDeleteAsset, 
  onAddTagClick, 
  onSaveTag, 
  onCancelAddTag, 
  onTagValueChange 
}) => {
  const getRiskColor = (risk) => {
    switch (risk) {
      case 'High':
        return 'bg-red-100 text-red-800 border-red-200';
      case 'Medium':
        return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'Low':
        return 'bg-green-100 text-green-800 border-green-200';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between mb-4">
        <span className={`px-3 py-1 rounded-full text-xs font-semibold border ${getRiskColor(asset.riskLevel)}`}>
          {asset.riskLevel} Risk
        </span>
        <button
          onClick={() => onDeleteAsset(asset.id, asset.hostname)}
          className="text-red-600 hover:text-red-800 transition-colors"
          title="Delete asset"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>

      <h3 className="text-lg font-semibold text-gray-900 mb-2 truncate" title={asset.hostname}>
        {asset.hostname}
      </h3>

      <p className="text-sm text-gray-600 mb-4 font-mono">
        {asset.ipAddress}
      </p>

      <div className="mb-4">
        <p className="text-xs text-gray-500 mb-2">Open Ports:</p>
        <div className="flex flex-wrap gap-2">
          {asset.ports.map((port) => (
            <span
              key={port}
              className="px-2 py-1 bg-gray-100 text-gray-700 rounded text-xs font-medium"
            >
              {port}
            </span>
          ))}
        </div>
      </div>

      <div>
        <div className="flex items-center justify-between mb-2">
          <p className="text-xs text-gray-500">Tags:</p>
          {editingTagForAsset !== asset.id && (
            <button
              onClick={() => onAddTagClick(asset.id)}
              className="text-xs text-blue-600 hover:text-blue-800 font-medium"
            >
              + Add Tag
            </button>
          )}
        </div>
        <div className="flex flex-wrap gap-2">
          {asset.tags.map((tag) => (
            <span
              key={tag}
              className="px-2 py-1 bg-blue-50 text-blue-700 rounded text-xs font-medium"
            >
              {tag}
            </span>
          ))}
          
          {editingTagForAsset === asset.id && (
            <div className="flex gap-2 w-full mt-2">
              <input
                type="text"
                value={newTagValue}
                onChange={(e) => onTagValueChange(e.target.value)}
                onKeyPress={(e) => {
                  if (e.key === 'Enter') {
                    onSaveTag(asset.id);
                  } else if (e.key === 'Escape') {
                    onCancelAddTag();
                  }
                }}
                placeholder="Enter tag name"
                className="flex-1 px-2 py-1 text-xs border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                autoFocus
              />
              <button
                onClick={() => onSaveTag(asset.id)}
                className="px-3 py-1 bg-blue-600 text-white text-xs rounded hover:bg-blue-700"
              >
                Save
              </button>
              <button
                onClick={onCancelAddTag}
                className="px-3 py-1 bg-gray-200 text-gray-700 text-xs rounded hover:bg-gray-300"
              >
                Cancel
              </button>
            </div>
          )}
        </div>
      </div>

      <div className="mt-4 pt-4 border-t border-gray-100">
        <p className="text-xs text-gray-400">
          Added {new Date(asset.createdAt).toLocaleDateString()}
        </p>
      </div>
    </div>
  );
};

export default AssetCard;
