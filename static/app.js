document.addEventListener('DOMContentLoaded', function() {
    // DOM Elements
    const uploadArea = document.getElementById('uploadArea');
    const fileInput = document.getElementById('fileInput');
    const browseBtn = document.getElementById('browseBtn');
    const fileInfo = document.getElementById('fileInfo');
    const fileName = document.getElementById('fileName');
    const migrateBtn = document.getElementById('migrateBtn');
    const progressBar = document.getElementById('progressBar');
    const progressFill = document.querySelector('.progress-fill');
    const resultsSection = document.getElementById('resultsSection');
    const resultMessage = document.getElementById('resultMessage');
    const downloadResultBtn = document.getElementById('downloadResultBtn');
    const downloadMigratedFileBtn = document.getElementById('downloadMigratedFileBtn');
    const sourceFormat = document.getElementById('sourceFormat');
    const targetFormat = document.getElementById('targetFormat');

    // Event Listeners
    browseBtn.addEventListener('click', () => {
        fileInput.click();
    });

    uploadArea.addEventListener('dragover', (e) => {
        e.preventDefault();
        uploadArea.style.borderColor = '#2980b9';
        uploadArea.style.backgroundColor = '#e3f2fd';
    });

    uploadArea.addEventListener('dragleave', () => {
        uploadArea.style.borderColor = '#3498db';
        uploadArea.style.backgroundColor = '#f8f9fa';
    });

    uploadArea.addEventListener('drop', (e) => {
        e.preventDefault();
        uploadArea.style.borderColor = '#3498db';
        uploadArea.style.backgroundColor = '#f8f9fa';
        
        if (e.dataTransfer.files.length) {
            handleFileSelection(e.dataTransfer.files[0]);
        }
    });

    fileInput.addEventListener('change', () => {
        if (fileInput.files.length) {
            handleFileSelection(fileInput.files[0]);
        }
    });

    migrateBtn.addEventListener('click', startMigration);
    downloadResultBtn.addEventListener('click', downloadResult);
    downloadMigratedFileBtn.addEventListener('click', downloadMigratedFile);

    // Functions
    function handleFileSelection(file) {
        fileName.textContent = file.name;
        fileInfo.classList.remove('hidden');
        migrateBtn.disabled = false;
        
        // Auto-detect source format based on file extension
        const ext = file.name.split('.').pop().toLowerCase();
        if (ext === 'sql') {
            sourceFormat.value = 'sql';
        } else if (ext === 'xlsx' || ext === 'xls') {
            sourceFormat.value = 'excel';
        } else if (ext === 'csv') {
            sourceFormat.value = 'csv';
        } else if (ext === 'json') {
            sourceFormat.value = 'json';
        } else if (ext === 'db' || ext === 'sqlite') {
            sourceFormat.value = 'sqlite';
        }
    }

    async function startMigration() {
        if (!fileInput.files.length) {
            alert('Please select a file first');
            return;
        }

        const file = fileInput.files[0];
        const sourceFmt = sourceFormat.value;
        const targetFmt = targetFormat.value;

        // Show progress bar
        progressBar.classList.remove('hidden');
        migrateBtn.disabled = true;

        // Create FormData object
        const formData = new FormData();
        formData.append('file', file);
        formData.append('sourceFormat', sourceFmt);
        formData.append('targetFormat', targetFmt);

        try {
            console.log('Sending request to /upload');
            
            // Send file to backend
            const response = await fetch('/upload', {
                method: 'POST',
                body: formData
            });

            console.log('Response status:', response.status);
            
            if (!response.ok) {
                const errorText = await response.text();
                console.error('Server error response:', errorText);
                throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`);
            }

            const result = await response.json();
            console.log('Server response:', result);
            
            // Simulate progress completion
            let progress = 0;
            const interval = setInterval(() => {
                progress += 10;
                progressFill.style.width = `${progress}%`;
                
                if (progress >= 100) {
                    clearInterval(interval);
                    completeMigration(result);
                }
            }, 100);
        } catch (error) {
            console.error('Error uploading file:', error);
            alert('Error uploading file: ' + error.message);
            progressBar.classList.add('hidden');
            migrateBtn.disabled = false;
        }
    }

    function completeMigration(result) {
        // Hide progress bar
        progressBar.classList.add('hidden');
        
        // Show results
        resultMessage.innerHTML = `
            <strong>Migration Complete!</strong><br>
            File: ${result.file}<br>
            Source Format: ${result.source.toUpperCase()}<br>
            Target Format: ${result.target.toUpperCase()}<br>
            Status: ${result.message}
        `;
        resultsSection.classList.remove('hidden');
        
        // Store migration info for download
        window.migrationResult = {
            fileName: result.file,
            sourceFormat: result.source,
            targetFormat: result.target,
            timestamp: new Date().toISOString(),
            content: result.message,
            migratedFile: result.migratedFile
        };
    }

    function downloadResult() {
        if (!window.migrationResult) {
            alert('No migration result available');
            return;
        }

        // Create a sample result file
        const content = `Migration Report
=================

File: ${window.migrationResult.fileName}
Source Format: ${window.migrationResult.sourceFormat.toUpperCase()}
Target Format: ${window.migrationResult.targetFormat.toUpperCase()}
Timestamp: ${window.migrationResult.timestamp}

Result:
${window.migrationResult.content}`;

        const blob = new Blob([content], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        
        const a = document.createElement('a');
        a.href = url;
        a.download = `migration-report-${Date.now()}.txt`;
        document.body.appendChild(a);
        a.click();
        
        // Clean up
        setTimeout(() => {
            document.body.removeChild(a);
            URL.revokeObjectURL(url);
        }, 100);
    }

    function downloadMigratedFile() {
        if (!window.migrationResult || !window.migrationResult.migratedFile) {
            alert('No migrated file available for download');
            return;
        }

        // Create download URL with the exact filename
        const downloadUrl = `/download-migrated?filename=${encodeURIComponent(window.migrationResult.migratedFile)}`;
        
        // Create a temporary link element
        const a = document.createElement('a');
        a.href = downloadUrl;
        
        // Determine file extension based on target format
        let extension = window.migrationResult.targetFormat;
        if (extension === 'excel') {
            extension = 'xlsx';
        }
        
        a.download = `migrated-${Date.now()}.${extension}`;
        document.body.appendChild(a);
        a.click();
        
        // Clean up
        setTimeout(() => {
            document.body.removeChild(a);
        }, 100);
    }

    // Initialize
    console.log('migrator Frontend initialized');
});