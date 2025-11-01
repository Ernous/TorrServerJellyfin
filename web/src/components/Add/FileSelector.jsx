import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import {
  Checkbox,
  FormControlLabel,
  FormGroup,
  FormHelperText,
  Typography,
  Box,
  Collapse,
  IconButton,
} from '@material-ui/core'
import { ExpandMore as ExpandMoreIcon, Folder as FolderIcon } from '@material-ui/icons'
import { styled } from '@material-ui/core/styles'

const FileSelectorContainer = styled(Box)(({ theme }) => ({
  marginTop: theme.spacing(2),
  marginBottom: theme.spacing(2),
}))

const ExpandButton = styled(IconButton)(({ theme, expanded }) => ({
  transform: expanded ? 'rotate(180deg)' : 'rotate(0deg)',
  transition: theme.transitions.create('transform', {
    duration: theme.transitions.duration.shortest,
  }),
}))

const FileItem = styled(FormControlLabel)(({ theme }) => ({
  marginLeft: theme.spacing(4),
  '& .MuiTypography-root': {
    fontSize: '0.875rem',
  },
}))

const DirectoryItem = styled(Box)(({ theme }) => ({
  marginLeft: theme.spacing(1),
  marginTop: theme.spacing(1),
}))

const DirectoryHeader = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  cursor: 'pointer',
  padding: theme.spacing(0.5),
  borderRadius: theme.shape.borderRadius,
  '&:hover': {
    backgroundColor: theme.palette.action.hover,
  },
}))

const DirectoryName = styled(Typography)(({ theme }) => ({
  marginLeft: theme.spacing(1),
  fontWeight: 500,
  fontSize: '0.875rem',
}))

export default function FileSelector({ torrentFiles, selectedFiles, setSelectedFiles }) {
  const { t } = useTranslation()
  const [expanded, setExpanded] = useState(false)
  const [selectAll, setSelectAll] = useState(true)
  const [expandedDirs, setExpandedDirs] = useState({})

  useEffect(() => {
    // Initialize with all files selected if selectedFiles is empty
    if (!selectedFiles || selectedFiles.length === 0) {
      const allFileIds = torrentFiles?.map(f => f.id) || []
      setSelectedFiles(allFileIds)
      setSelectAll(true)
    } else {
      setSelectAll(selectedFiles.length === torrentFiles?.length)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [torrentFiles])

  const handleSelectAll = event => {
    const { checked } = event.target
    setSelectAll(checked)
    if (checked) {
      const allFileIds = torrentFiles?.map(f => f.id) || []
      setSelectedFiles(allFileIds)
    } else {
      setSelectedFiles([])
    }
  }

  const handleFileToggle = fileId => {
    const currentIndex = selectedFiles.indexOf(fileId)
    const newSelected = [...selectedFiles]

    if (currentIndex === -1) {
      newSelected.push(fileId)
    } else {
      newSelected.splice(currentIndex, 1)
    }

    setSelectedFiles(newSelected)
    setSelectAll(newSelected.length === torrentFiles?.length)
  }

  const handleDirToggle = dirPath => {
    setExpandedDirs(prev => ({
      ...prev,
      [dirPath]: !prev[dirPath]
    }))
  }

  // Группировка файлов по директориям
  const groupFilesByDirectory = files => {
    const grouped = {}
    
    files.forEach(file => {
      const pathParts = file.path.split('/')
      if (pathParts.length === 1) {
        // Файл в корне
        if (!grouped['']) grouped[''] = []
        grouped[''].push(file)
      } else {
        // Файл в директории
        const dirPath = pathParts.slice(0, -1).join('/')
        if (!grouped[dirPath]) grouped[dirPath] = []
        grouped[dirPath].push(file)
      }
    })
    
    return grouped
  }

  const handleDirSelectAll = (dirFiles, select) => {
    const dirFileIds = dirFiles.map(f => f.id)
    const newSelected = [...selectedFiles]
    
    if (select) {
      // Добавляем все файлы директории
      dirFileIds.forEach(id => {
        if (!newSelected.includes(id)) {
          newSelected.push(id)
        }
      })
    } else {
      // Убираем все файлы директории
      dirFileIds.forEach(id => {
        const index = newSelected.indexOf(id)
        if (index !== -1) {
          newSelected.splice(index, 1)
        }
      })
    }
    
    setSelectedFiles(newSelected)
    setSelectAll(newSelected.length === torrentFiles?.length)
  }

  if (!torrentFiles || torrentFiles.length === 0) {
    return null
  }

  const videoFiles = torrentFiles.filter(f => {
    const ext = f.path.toLowerCase().split('.').pop()
    return ['mkv', 'mp4', 'avi', 'mov', 'wmv', 'flv', 'webm', 'm4v'].includes(ext)
  })

  if (videoFiles.length === 0) {
    return null
  }

  const groupedFiles = groupFilesByDirectory(videoFiles)
  const directories = Object.keys(groupedFiles).sort()

  return (
    <FileSelectorContainer>
      <Box display='flex' alignItems='center' justifyContent='space-between'>
        <Typography variant='subtitle2'>{t('AddDialog.SelectFiles', 'Select files for .strm creation')}</Typography>
        <ExpandButton
          expanded={expanded ? 1 : 0}
          onClick={() => setExpanded(!expanded)}
          aria-expanded={expanded}
          aria-label='show more'
          size='small'
        >
          <ExpandMoreIcon />
        </ExpandButton>
      </Box>

      <FormHelperText>
        {selectedFiles.length} / {videoFiles.length} {t('AddDialog.FilesSelected', 'files selected')}
      </FormHelperText>

      <Collapse in={expanded} timeout='auto' unmountOnExit>
        <FormGroup>
          <FormControlLabel
            control={
              <Checkbox
                checked={selectAll}
                onChange={handleSelectAll}
                indeterminate={selectedFiles.length > 0 && selectedFiles.length < videoFiles.length}
                color='primary'
              />
            }
            label={
              <Typography variant='body2'>
                <strong>{t('SelectAll', 'Select All')}</strong>
              </Typography>
            }
          />
          
          {directories.map(dirPath => {
            const dirFiles = groupedFiles[dirPath]
            const selectedInDir = dirFiles.filter(f => selectedFiles.includes(f.id)).length
            const isExpanded = expandedDirs[dirPath] !== false // По умолчанию раскрыты
            const isDirSelected = selectedInDir === dirFiles.length
            const isDirIndeterminate = selectedInDir > 0 && selectedInDir < dirFiles.length
            
            return (
              <DirectoryItem key={dirPath || 'root'}>
                <DirectoryHeader onClick={() => handleDirToggle(dirPath)}>
                  <Checkbox
                    checked={isDirSelected}
                    indeterminate={isDirIndeterminate}
                    onChange={e => handleDirSelectAll(dirFiles, e.target.checked)}
                    onClick={e => e.stopPropagation()}
                    color='primary'
                    size='small'
                  />
                  <FolderIcon fontSize='small' color='action' />
                  <DirectoryName variant='body2'>
                    {dirPath || t('RootDirectory', 'Root')} ({selectedInDir}/{dirFiles.length})
                  </DirectoryName>
                  <ExpandButton
                    expanded={isExpanded ? 1 : 0}
                    size='small'
                  >
                    <ExpandMoreIcon />
                  </ExpandButton>
                </DirectoryHeader>
                
                <Collapse in={isExpanded} timeout='auto'>
                  <FormGroup>
                    {dirFiles.map(file => {
                      const fileName = file.path.split('/').pop()
                      return (
                        <FileItem
                          key={file.id}
                          control={
                            <Checkbox
                              checked={selectedFiles.indexOf(file.id) !== -1}
                              onChange={() => handleFileToggle(file.id)}
                              color='primary'
                              size='small'
                            />
                          }
                          label={fileName}
                        />
                      )
                    })}
                  </FormGroup>
                </Collapse>
              </DirectoryItem>
            )
          })}
        </FormGroup>
      </Collapse>
    </FileSelectorContainer>
  )
}
