import React, { useState, useEffect, useCallback } from 'react';
import CreateAssetModal from './CreateAssetModal';
import Checkbox from './Checkbox';
import AssetCard from './AssetCard';
import Pagination from './Pagination';

const AssetList = () => {
  const [assets, setAssets] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [riskFilter, setRiskFilter] = useState(['High', 'Medium', 'Low']);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [editingTagForAsset, setEditingTagForAsset] = useState(null);
  const [newTagValue, setNewTagValue] = useState('');
  const [pagination, setPagination] = useState({ total: 0, page: 1, page_size: 10 });
  const itemsPerPage = 6;

  // Fetch assets from API
  const fetchAssets = useCallback(async (page = currentPage) => {
    setLoading(true);
    setError(null);

    try {
      const params = new URLSearchParams();
      
      // Add search parameter
      if (searchTerm) {
        params.append('search', searchTerm);
      }

      // Add risk level parameter (comma-separated)
      if (riskFilter.length > 0) {
        params.append('risk_level', riskFilter.join(','));
      }

      // Add pagination parameters
      params.append('page', page);
      params.append('page_size', itemsPerPage);

      const url = `http://localhost:8080/api/v1/assets${params.toString() ? '?' + params.toString() : ''}`;
      const response = await fetch(url);

      if (!response.ok) {
        throw new Error('Failed to fetch assets');
      }

      const data = await response.json();

      // Transform API response to match frontend format
      const transformedAssets = (data.assets || []).map(asset => ({
        id: asset.id,
        ipAddress: asset.ip_address,
        hostname: asset.hostname,
        riskLevel: asset.risk_level,
        ports: (asset.ports || []).map(p => p.port_number || p),
        tags: (asset.tags || []).map(t => t.tag_name || t),
        createdAt: asset.created_at
      }));

      setAssets(transformedAssets);
      setPagination(data.pagination || { total: 0, page: 1, page_size: itemsPerPage });
    } catch (err) {
      setError(err.message);
      setAssets([]);
    } finally {
      setLoading(false);
    }
  }, [searchTerm, riskFilter, itemsPerPage]);

  // Debounce search term and fetch assets
  useEffect(() => {
    // Debounce search: wait 500ms after user stops typing
    const debounceTimer = setTimeout(() => {
      fetchAssets(1); // Reset to page 1 when filters change
      setCurrentPage(1);
    }, 500);

    // Cleanup: cancel the timer if searchTerm or riskFilter changes before 500ms
    return () => clearTimeout(debounceTimer);
  }, [searchTerm, riskFilter, fetchAssets]);

  // Fetch assets when page changes
  useEffect(() => {
    fetchAssets(currentPage);
  }, [currentPage, fetchAssets]);

  const handleAssetCreated = (newAsset) => {
    // Refresh the asset list after creation
    fetchAssets();
    setCurrentPage(1);
  };

  const handleAddTagClick = (assetId) => {
    setEditingTagForAsset(assetId);
    setNewTagValue('');
  };

  const handleCancelAddTag = () => {
    setEditingTagForAsset(null);
    setNewTagValue('');
  };

  const handleSaveTag = async (assetId) => {
    if (!newTagValue.trim()) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/api/v1/assets/${assetId}/tags`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          tag_name: newTagValue.trim(),
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to add tag');
      }

      // Refresh assets to show new tag
      await fetchAssets();
      setEditingTagForAsset(null);
      setNewTagValue('');
    } catch (err) {
      console.error('Error adding tag:', err);
      setError('Failed to add tag');
    }
  };

  const handleDeleteAsset = async (assetId, hostname) => {
    if (!window.confirm(`Are you sure you want to delete "${hostname}"? This action cannot be undone.`)) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/api/v1/assets/${assetId}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error('Failed to delete asset');
      }

      // Refresh assets list
      await fetchAssets(currentPage);
    } catch (err) {
      console.error('Error deleting asset:', err);
      setError('Failed to delete asset');
    }
  };

  // Calculate pagination from API response
  const totalPages = Math.ceil((pagination.total || 0) / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = Math.min(startIndex + itemsPerPage, pagination.total || 0);

  // Reset to page 1 when filters change
  const handleSearchChange = (value) => {
    setSearchTerm(value);
    setCurrentPage(1);
  };

  const handleRiskFilterChange = (level) => {
    setRiskFilter(prev => {
      if (prev.includes(level)) {
        // Remove if already selected
        return prev.filter(l => l !== level);
      } else {
        // Add if not selected
        return [...prev, level];
      }
    });
    setCurrentPage(1);
  };

  const toggleAllRiskLevels = () => {
    if (riskFilter.length === 3) {
      // If all selected, deselect all
      setRiskFilter([]);
    } else {
      // Otherwise, select all
      setRiskFilter(['High', 'Medium', 'Low']);
    }
    setCurrentPage(1);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Asset Management</h1>
          <p className="mt-2 text-sm text-gray-600">
            Monitor and manage your network assets
          </p>
        </div>

        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
          <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
            <div className="flex flex-col sm:flex-row gap-4 flex-1 w-full sm:w-auto">
              <div className="flex-1 min-w-0">
                <input
                  type="text"
                  placeholder="Search by hostname or IP..."
                  value={searchTerm}
                  onChange={(e) => handleSearchChange(e.target.value)}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                />
              </div>

              <div className="sm:w-auto">
                <div className="bg-white border border-gray-300 rounded-lg p-3">
                  <div className="flex items-center gap-4">
                    <span className="text-sm font-medium text-gray-700">Risk Level:</span>
                    
                    <Checkbox
                      checked={riskFilter.length === 3}
                      onChange={toggleAllRiskLevels}
                      label="All"
                      colorClass="gray"
                    />

                    <Checkbox
                      checked={riskFilter.includes('High')}
                      onChange={() => handleRiskFilterChange('High')}
                      label="High"
                      colorClass="red"
                    />

                    <Checkbox
                      checked={riskFilter.includes('Medium')}
                      onChange={() => handleRiskFilterChange('Medium')}
                      label="Medium"
                      colorClass="yellow"
                    />

                    <Checkbox
                      checked={riskFilter.includes('Low')}
                      onChange={() => handleRiskFilterChange('Low')}
                      label="Low"
                      colorClass="green"
                    />
                  </div>
                </div>
              </div>
            </div>

            <button
              onClick={() => setIsModalOpen(true)}
              className="w-full sm:w-auto px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium shadow-sm"
            >
              + Create Asset
            </button>
          </div>

          <div className="mt-4 text-sm text-gray-600">
            {loading ? (
              'Loading assets...'
            ) : (
              <>
                Showing {pagination.total > 0 ? startIndex + 1 : 0}-{Math.min(endIndex, pagination.total)} of {pagination.total} assets
              </>
            )}
          </div>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
            <p className="text-sm text-red-800">Error: {error}</p>
          </div>
        )}

        {loading && (
          <div className="flex justify-center items-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
          </div>
        )}

        {!loading && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {assets.map((asset) => (
              <AssetCard
                key={asset.id}
                asset={asset}
                editingTagForAsset={editingTagForAsset}
                newTagValue={newTagValue}
                onDeleteAsset={handleDeleteAsset}
                onAddTagClick={handleAddTagClick}
                onSaveTag={handleSaveTag}
                onCancelAddTag={handleCancelAddTag}
                onTagValueChange={setNewTagValue}
              />
            ))}
          </div>
        )}

        {!loading && assets.length === 0 && (
          <div className="text-center py-12">
            <h3 className="text-lg font-medium text-gray-900 mb-2">No assets found</h3>
            <p className="text-gray-500">Try adjusting your search or filters</p>
          </div>
        )}

        {pagination.total > 0 && (
          <Pagination
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={setCurrentPage}
          />
        )}
      </div>

      <CreateAssetModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onAssetCreated={handleAssetCreated}
      />
    </div>
  );
};

export default AssetList;
